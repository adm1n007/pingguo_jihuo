package inputRecord

import (
    . "ml/strings"
    "ituneslib"
)

type ElementPosition struct {
    Name    String
    Left    int
    Right   int
    Top     int
    Bottom  int
}

type ElementLayout []*ElementPosition

func LayoutForCountry(country ituneslib.CountryID) ElementLayout {
    return Layout_PC[country]
}

func LayoutForCountryiOS(country ituneslib.CountryID) ElementLayout {
    return Layout_iOS[country]
}

var Layout_PC = map[ituneslib.CountryID]ElementLayout{
    ituneslib.CountryID_China: ElementLayout{
        { "codeRedemptionField",    120, 420,  370, 390 },
        { "lastFirstName",          120, 260,  480, 500 },
        { "firstName",              280, 420,  480, 500 },
        { "street1",                120, 420,  515, 535 },
        { "street2",                120, 420,  550, 570 },
        { "street3",                120, 420,  585, 605 },
        { "city",                   120, 240,  620, 640 },
        { "postalcode",             120, 215,  655, 675 },
        { "state",                  230, 350,  655, 675 },
        { "phone1Number",           120, 240,  690, 710 },
        { "continue",               920, 1020, 850, 870 },
    },

    /*
        [输入代码]

        [称呼]
        [姓    氏] [名    字]
        [县市]
        [街                   道]
        [公寓、套房、大楼         ]
        [邮递区号   ] Taiwan
        [区域号码][电       话]

        [codeRedemptionField]

        [Title]
        [First Name] [Last Name]
        [Cyty]
        [Street1               ]
        [Street2               ]
        [postalcode] Taiwan
        [AreaCode][Phone  ]
    */
    ituneslib.CountryID_Taiwan: ElementLayout{
        { "codeRedemptionField",    186, 496,  352, 377 },


        { "salutation",             186, 261,  465, 490 },

        { "lastFirstName",          186, 336,  500, 525 },
        { "firstName",              347, 497,  500, 525 },

        { "citypopup",              186, 261,  535, 560 },

        { "street1",                186, 497,  570, 595 },
        { "street2",                186, 497,  605, 630 },

        { "postalcode",             186, 287,  640, 665 },
        { "phone1AreaCode",         186, 274,  675, 700 },
        { "phone1Number",           274, 399,  675, 700 },

        { "continue",               927, 1095, 838, 861 },
    },

    ituneslib.CountryID_Japan: ElementLayout{
        { "codeRedemptionField",    186, 496,   368, 393 },
        { "lastFirstName",          186, 336,   481, 506 },
        { "firstName",              347, 497,   481, 506 },
        { "phoneticLastName",       186, 336,   516, 541 },
        { "phoneticFirstName",      347, 497,   516, 541 },
        { "postalcode",             186, 287,   551, 576 },
        { "state",                  186, 266,   586, 611 },
        { "city",                   276, 404,   588, 613 },
        { "street1",                186, 497,   623, 648 },
        { "street2",                186, 497,   658, 683 },
        { "phone1AreaCode",         186, 274,   693, 718 },
        { "phone1Number",           274, 399,   693, 718 },
        { "continue",               933, 1095,  852, 871 },
    },

    ituneslib.CountryID_NewZealand: ElementLayout{
        { "codeRedemptionField",    120, 420,   370, 390 },
        { "salutation",             120, 150,   480, 500 },
        { "firstName",              120, 260,   515, 535 },
        { "lastName",               280, 420,   515, 535 },
        { "street1",                120, 420,   550, 570 },
        { "street2",                120, 420,   585, 605 },
        { "suburb",                 120, 240,   620, 640 },
        { "postalcode",             120, 200,   655, 675 },
        { "city",                   240, 300,   655, 675 },
        { "phone1AreaCode",         120, 180,   690, 710 },
        { "phone1Number",           200, 250,   690, 710 },
        { "continue",               920, 1020,  850, 870 },
    },
}

var Layout_iOS = map[ituneslib.CountryID]ElementLayout{
    ituneslib.CountryID_Taiwan: ElementLayout{
        { "codeRedemptionField",        137,    305,   353, 378 },
        { "salutationField",            137,    305,   456, 477 },
        { "lastNameField",              137,    305,   499, 524 },
        { "firstNameField",             137,    305,   544, 569 },
        { "cityField",                  137,    305,   591, 612 },
        { "street1Field",               137,    305,   634, 659 },
        { "street2Field",               137,    305,   679, 704 },
        { "postalCodeField",            137,    305,   724, 749 },
        { "phoneAreaCodeField",         137,    192,   800, 843 },
        { "phoneNumberField",           192,    305,   809, 834 },
        { "hiddenBottomRightButtonId",  -2248,  -2048,  1600, 1632 },
    },
}
