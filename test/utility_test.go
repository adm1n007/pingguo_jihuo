package test

import (
    "testing"
    . "fmt"
    . "ml/strings"

    "encoding/json"
)

var q String = `
  var questionsDictionary = {"securityQuestions":[{"id":130,"question":"你少年时代最好的朋友叫什么名字？"},{"id":131,"question":"你的第一个宠物叫什么名字？"},{"id":132,"question":"你学会做的第一道菜是什么？"},{"id":133,"question":"你第一次去电影院看的是哪一部电影？"},{"id":134,"question":"你第一次坐飞机是去哪里？"},{"id":135,"question":"你上小学时最喜欢的老师姓什么？"},{"id":136,"question":"你的理想工作是什么？"},{"id":137,"question":"你小时候最喜欢哪一本书？"},{"id":138,"question":"你拥有的第一辆车是什么型号？"},{"id":139,"question":"你童年时代的绰号是什么？"},{"id":140,"question":"你在学生时代最喜欢哪个电影明星或角色？"},{"id":141,"question":"你在学生时代最喜欢哪个歌手或乐队？"},{"id":142,"question":"你的父母是在哪里认识的？"},{"id":143,"question":"你的第一个上司叫什么名字？"},{"id":144,"question":"您从小长大的那条街叫什么？"},{"id":145,"question":"你去过的第一个海滨浴场是哪一个？"},{"id":146,"question":"你购买的第一张专辑是什么？"},{"id":147,"question":"您最喜欢哪个球队？"}],"securityQuestionsPerPage":6};
`

func TestUtility(t *testing.T) {

    type Question struct {
        Id          int
        Question    string
    }

    type Questions struct {
        SecurityQuestions           []Question
        SecurityQuestionsPerPage    int
    }

    questions := Questions{}

    q = "{" + q.Split("{", 1)[1].RSplit(";", 1)[0]
    json.Unmarshal(q.Encode(CP_UTF8), &questions)

    // Println(q)
    Println(questions)
}
