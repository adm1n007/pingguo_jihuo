package register
import (
    . "fmt"
    . "ml/trace"
    . "ml/dict"
    . "ml/array"
    . "ml/strings"

    "time"
    "encoding/base64"
    urllib "net/url"

    "github.com/PuerkitoBio/goquery"
    "ml/random"
    "ituneslib"

    "utility"
)

func (self *AppleIdRegister) getProvideAppleIdDetailsInfo(doc *goquery.Document, elementIds Array, callback fillFormCallback) (Dict, elementDict) {
    self.debug("getProvideAppleIdDetailsInfo")

    data := getAllElement(doc.Selection, elementIds)
    if len(data) != len(elementIds) {
        Raise(NewGenericError("len(data) != len(elementIds), resp incorrect"))
    }

    data["continue"] = doc.Find("input.continue[type=submit]")

    var birthYear int
    var birthMonth int
    var birthDay int

    if self.account.InitFromDatabase {
        birthday := String(self.account.Birth).Split("-")
        birthYear = birthday[0].ToInt()
        birthMonth = birthday[1].ToInt()
        birthDay = birthday[2].ToInt()
    } else {
        birthYear = random.IntRange(1970, 2001)
        birthMonth = random.IntRange(1, 13)
        birthDay = random.IntRange(1, 27)
    }

    _ = birthMonth
    _ = birthDay

    form := Dict{
        //
        // apple id and password
        //

        data["emailAddress"].Attr2("name")      : self.account.UserName,
        data["pass1"].Attr2("name")             : self.account.ApplePassword,
        data["pass2"].Attr2("name")             : self.account.ApplePassword,

        //
        // apple id and password
        //

        data["birthYear"].Attr2("name")         : birthYear,
        data["birthMonthPopup"].Attr2("name")   : selectOption(data["birthMonthPopup"]),
        data["birthDayPopup"].Attr2("name")     : selectOption(data["birthDayPopup"], Dict{"maxIndex": 25}),
        // data["birthMonthPopup"].Attr2("name")   : data["birthMonthPopup"].Find(Sprintf("option:contains('%d')", birthMonth)).Attr2("value"),
        // data["birthDayPopup"].Attr2("name")     : data["birthDayPopup"].Find(Sprintf("option:contains('%d')", birthDay)).Attr2("value"),

        //
        // subscribe
        //

        data["newsletter"].Attr2("name")        : getInputValue(data["newsletter"]),
        data["marketing"].Attr2("name")         : getInputValue(data["marketing"]),

        //
        // continue
        //

        data["continue"].Attr2("name")          : data["continue"].Attr2("value"),

        //
        // signature
        //

        data["pageUUID"].Attr2("name")          : data["pageUUID"].Attr2("value"),
        data["machineGUID"].Attr2("name")       : self.machineGUID,
        data["signature"].Attr2("name")         : base64.StdEncoding.EncodeToString(self.sapSession.SignData((String(data["pageUUID"].Attr2("value")) + self.machineGUID + String(self.account.UserName)).Encode(CP_UTF8))),
    }

    year := form[data["birthYear"].Attr2("name")]
    monthValue := String(form[data["birthMonthPopup"].Attr2("name")].(string))
    dayValue := form[data["birthDayPopup"].Attr2("name")]

    month := String(data["birthMonthPopup"].Find(Sprintf(`option[value="%v"]`, monthValue)).Text())
    day := data["birthDayPopup"].Find(Sprintf(`option[value="%v"]`, dayValue)).Text()

    if month.IsDigit() == false {
        for m := time.January; m <= time.December; m++ {
            if String(m.String()).ToLower() == month.ToLower() {
                month = String(Sprintf("%d", m))
                break
            }
        }
    }

    self.account.Birth = Sprintf("%v-%v-%v", year, month, day)

    callback(form, data)

    return form, data
}

func (self *AppleIdRegister) stepProvideAppleIdDetails(refer *urllib.URL, next string, form Dict) bool {
    self.info(`stepProvideAppleIdDetails`)

    resp := self.post(
                buildNext(refer, next),
                Dict{
                    "headers": Dict{
                        "Referer"           : refer.String(),
                        "Origin"            : buildOrigin(refer),
                        "Accept"            : "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
                        "Accept-Encoding"   : "gzip, deflate",
                        "Content-Type"      : "application/x-www-form-urlencoded",
                        `Connection`        : `keep-alive`,
                    },
                    "body": form,
                },
            )

    doc := utility.ParseHTML(resp.Text())
    self.checkPageError(doc)
    self.pingUrl(doc, resp)

    switch self.account.Country {
        case ituneslib.CountryID_Vietnam:
            form = self.stepProvideAppleIdDetails_Vietnam(doc)

        default:
            form = self.stepProvideAppleIdDetails_Default(doc)
    }

    return self.stepSelectPayEase(
        resp.Request.URL,
        doc.Find("form[method=post]").Attr2("action"),
        form,
    )
}

func (self *AppleIdRegister) stepProvideAppleIdDetails_Default(doc *goquery.Document) Dict {
    self.debug("stepProvideAppleIdDetails_Default")

    elementIds := Array{
        "emailAddress",
        "pass1",
        "pass2",
        "question1",
        "question2",
        "question3",
        "answer1",
        "answer2",
        "answer3",
        "recoveryEmail",
        "birthYear",
        "birthMonthPopup",
        "birthDayPopup",
        "newsletter",
        "marketing",
        "pageUUID",
        "machineGUID",
        "signature",
    }

    form, _ := self.getProvideAppleIdDetailsInfo(doc, elementIds, func (form Dict, elements elementDict) {
        questions := self.parseQuestionsDictionary(doc)

        var questionIndex1, questionIndex2, questionIndex3 string
        var questionText1, questionText2, questionText3 string
        var answer1, answer2, answer3 string

        if self.account.InitFromDatabase == false {
            questionIndex1 = selectOption(elements["question1"])
            questionIndex2 = selectOption(elements["question2"])
            questionIndex3 = selectOption(elements["question3"])

            questionText1 = elements["question1"].Find(Sprintf(`option[value="%v"]`, questionIndex1)).Text()
            questionText2 = elements["question2"].Find(Sprintf(`option[value="%v"]`, questionIndex2)).Text()
            questionText3 = elements["question3"].Find(Sprintf(`option[value="%v"]`, questionIndex3)).Text()

            answer1 = utility.GenerateAnswer().String()
            answer2 = utility.GenerateAnswer().String()
            answer3 = utility.GenerateAnswer().String()

        } else {
            getQuestionIndex := func (questionIndex string, questionText string) string {
                q := elements[questionIndex].Find(Sprintf(`option:contains('%s')`, questionText))
                if q != nil && q.Length() != 0 {
                    return q.Attr2("value")
                }

                self.debug("cant find question %q in %q", questionText, questionIndex)
                Raisef("cant find question %q in %q", questionText, questionIndex)

                return ""
            }

            questionText1 = questions.QuestionsTextTable[String(self.account.Question1).ToInt()]
            questionText2 = questions.QuestionsTextTable[String(self.account.Question2).ToInt()]
            questionText3 = questions.QuestionsTextTable[String(self.account.Question3).ToInt()]

            questionIndex1 = getQuestionIndex("question1", questionText1)
            questionIndex2 = getQuestionIndex("question2", questionText2)
            questionIndex3 = getQuestionIndex("question3", questionText3)

            answer1 = self.account.Answer1
            answer2 = self.account.Answer2
            answer3 = self.account.Answer3
        }

        e := Dict{
            //
            // recoveryEmail
            //

            elements["recoveryEmail"].Attr2("name") : "",


            //
            // 3 questions and answers
            //

            elements["question1"].Attr2("name")     : questionIndex1,
            elements["answer1"].Attr2("name")       : answer1,
            elements["question2"].Attr2("name")     : questionIndex2,
            elements["answer2"].Attr2("name")       : answer2,
            elements["question3"].Attr2("name")     : questionIndex3,
            elements["answer3"].Attr2("name")       : answer3,
        }

        mergeDict(form, e)

        self.debug("questionText1 = %v", questionText1)
        self.debug("questionText2 = %v", questionText2)
        self.debug("questionText3 = %v", questionText3)
        self.debug("questionIndex1 = %v", questionIndex1)
        self.debug("questionIndex2 = %v", questionIndex2)
        self.debug("questionIndex3 = %v", questionIndex3)

        self.account.RecoveryEmail  = Sprintf("%v", form[elements["recoveryEmail"].Attr2("name")])
        self.account.Question1      = Sprintf("%v", questions.QuestionsIdTable[questionText1])
        self.account.Answer1        = answer1
        self.account.QuestionText1  = questionText1
        self.account.Question2      = Sprintf("%v", questions.QuestionsIdTable[questionText2])
        self.account.Answer2        = answer2
        self.account.QuestionText2  = questionText2
        self.account.Question3      = Sprintf("%v", questions.QuestionsIdTable[questionText3])
        self.account.Answer3        = answer3
        self.account.QuestionText3  = questionText3

        self.debug("%v", form)
        self.debug("%+v", self.account)
    })

    return form
}

func (self *AppleIdRegister) stepProvideAppleIdDetails_Vietnam(doc *goquery.Document) Dict {
    self.debug("stepProvideAppleIdDetails_Default")

    elementIds := Array{
        "emailAddress",
        "pass1",
        "pass2",
        "question",
        "answer",
        "birthYear",
        "birthMonthPopup",
        "birthDayPopup",
        "pageUUID",
        "machineGUID",
        "signature",
    }

    form, _ := self.getProvideAppleIdDetailsInfo(doc, elementIds, func (form Dict, elements elementDict) {
        question := utility.GenerateRandomString()[:random.IntRange(8, 12)]
        answer := utility.GenerateRandomString()[:random.IntRange(8, 12)]

        e := Dict{
            elements["question"].Attr2("name")  : question,
            elements["answer"].Attr2("name")    : answer,
        }

        mergeDict(form, e)

        self.account.Question1      = "-1"
        self.account.Answer1        = Sprintf("%v", form[elements["answer"].Attr2("name")])
        self.account.QuestionText1  = question.String()
    })

    return form
}
