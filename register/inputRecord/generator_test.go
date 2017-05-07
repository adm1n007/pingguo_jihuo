package inputRecord

import (
    . "ml/strings"
    . "ml/dict"
    "testing"
    "fmt"
    "github.com/PuerkitoBio/goquery"
    "ituneslib"
    "ml/html"
    "../../htmlutils"
)

func TestGenerate(t *testing.T) {
    doc := html.Parse(String(resp))

    g := New(LayoutForCountry(ituneslib.CountryID_China))

    elems := []*goquery.Selection{}

    htmlutils.FindInputs(doc.Selection).Each(func (i int, s *goquery.Selection) {
        elems = append(elems, s)
    })

    form := Dict{
        "longName": "宦灵,殊",
        "2.0.1.1.3.0.7.11.3.1.0.5.23.5.7.13.3": "68区609宿舍",
        "2.0.1.1.3.0.7.11.3.1.0.5.21.1.3.1.3.2.3": "",
        "ndpd-ipr": "ncip,0,56e7bbf3,1,1;st,0,mobile-phone,0,codeRedemptionField,0,lastFirstName,0,firstName,0,street1,0,street2,0,street3,0,city,0,postalcode,0,phone1Number,0;mm,8c,131,1d7,firstName;kk,859,0,codeRedemptionField;ff,0,codeRedemptionField;mc,6d,166,162,codeRedemptionField;fb,d42,codeRedemptionField;kk,1,0,lastFirstName;ff,1,lastFirstName;kd,68e2;ts,0,7f78;kd,35;kd,251;kd,2c3;kd,19b;mm,5d2,fd,284,state;kd,120;kd,97;fb,49d,lastFirstName;kk,1,0,firstName;ff,1,firstName;mc,55,156,1d5,firstName;kd,289c;ts,0,bb75;kd,d7;kd,1c1;fb,481,firstName;mm,178,117,1fc,street1;kk,0,0,street1;ff,0,street1;mc,62,9b,1f4,street1;kd,41b1;ts,0,10619;kd,186;kd,e7;kd,246;mm,3ed,d2,234,street3;kd,d8;kd,1f9;kd,28b;kd,186;kd,90;kd,1c2;kd,1e1;kd,220;kd,161;kd,c3;kd,1fa;kd,188;kd,134;kd,24c;kd,1fe;kd,2ab;kd,270;kd,12f;kd,da;mm,449,111,210,street2;kd,5f;kd,2cc;kd,2fc;kd,91;kd,187;kd,236;kd,106;fb,4d3,street1;ts,0,142cd;kk,0,0,street2;ff,1,street2;mc,5e,14f,215,street2;kd,1787;kd,10b;kd,129;kd,bd;kd,158;kd,1a5;mm,395,a2,266,city;kd,28e;kd,209;kd,13f;kd,248;kd,1b8;kd,115;kd,7a;kd,197;kd,c7;fb,582,street2;kk,0,0,street3;ff,1,street3;mc,4d,18a,241,street3;kd,6fe6;ts,0,1e7af;kd,123;mm,52d,fa,1cf,lastFirstName;kd,270;kd,258;kd,168;kd,87;kd,cf;kd,23f;kd,75;kd,2e6;kd,35;kd,15d;kd,28d;kd,197;kd,157;kd,1b8;kd,39;kd,1d0;kd,fd;mm,e8,af,280,postalcode;kd,254;kd,93;fb,57d,street3;kk,1,0,city;ff,1,city;mc,59,96,256,city;kd,3109;ts,0,2419a;kd,1b0;kd,178;kd,d3;kd,206;kd,2c2;kd,5f;kd,2c8;kd,224;kd,1d8;kd,253;mm,1da,8c,255,city;fb,3f5,city;kk,0,0,postalcode;ff,1,postalcode;kd,5d84;ts,0,2b727;kd,1cb;kd,1fd;kd,2e2;kd,243;kd,1f6;kd,3b;fb,462,postalcode;mm,267,7e,29b,phone1Number;kk,1,0,phone1Number;ff,0,phone1Number;mc,53,96,29d,phone1Number;kd,7a7b;ts,0,343dd;mm,28f,19f,1fe,street1;kd,c4;kd,28d;kd,26b;kd,9d;kd,da;kd,9e;kd,af;kd,1ff;kd,28a;kd,27d;kd,2bc;fb,a8a,phone1Number;mc,fb,3b8,33e,2.0.1.1.3.0.7.11.3.9.1;fs,5,0,0,;",
        "ndpd-fm": "1.w-855182.1.2.Js47oLxGDkYzWvEmaQmcbA,,.TVfBnDiCCNWABb_nfd9oIC26C-UJAXKYRQ0RPJu7MYzUFoTcT6BegGSEFf4SB6o5gEIoBSRMER5xgbUv1Y81SE2L4GiQLd8UBKOaSKzE4FxYOPwE3131UUUENSRyQRGalQp1dB3Y4hPp2rXECEADHHEJl1Nx8WIXAVLVxtVY1luwQzJSTtG8EJ6voF3KB8ds1V3hfjpyqFfXMfISizTD2AnsEvm3_tbuDu9KVZiTCNrNcNA_4lG4K3tKjg8dmSPK6DgKAPwoXoJxsGKFzLc3s8mgPm78-Rganles7ANJqpRb09sP5nTvS1KDAHjQ___2_uIRWVP7OWFxTStwda33yg,,",
        "machineGUID": "7a3135231bdf",
        "sp": "",
        "ndpd-w": "1.w-855182.1.2.kltcFW-yNWnb4D6mF1hL9Q,,.92467fsMwqe8SDT7gA92nh8qth2Lf-BwWXjobv1kU_xVRWbKCaDDNimWi5idBFRm3cbH4Uc4Vt_aLTAo027OVTHyyfwd_OdW2bLcR-j4RsGdaq0I8r_scrbMq1DtPklJ_qXnVTqJsLTGVJM57QNSI4qpVZzxWHI76J4O0OQxLL7eVphe9diB_6_nZmFu5hlu7sGYl7LvbzCA6KeeLCia9VkDWeufqG4OAErhV1H1mhIkdBPP9EGc6uoOvVuYejWgoyY--2k7qOUhaJKtYz_9xT1lumWebnOxLmRMGii2Qi7WE-q9jQwAWS-kPDDQiya2ILrb3ktaVpZDy_3-c6Sq3l6lrIF5O7i2oILmV5V0pLk9pKevCAbGzp7EXH-BUyLa",
        "ndpd-wk": "p",
        "captchaMode": "VIDEO",
        "2.0.1.1.3.0.7.11.3.9.1": "创建 Apple ID",
        "credit-card-type": "",
        "2.0.1.1.3.0.7.11.3.1.0.5.23.5.7.13.1": "炮炽稼胡同99大院",
        "2.0.1.1.3.0.7.11.3.1.0.5.23.5.7.13.5": "伪佣睁謇路",
        "res": "",
        "ndpd-wkr": 357529,
        "mzPageUUID": "oGPXeSwTStYSP6YN65TWzXxTZ/0=",
        "ndpd-s": "1.w-855182.1.2.EYXIj_nUPKdSi1uv4ypUnQ,,.gtnwcolCqH_Gal6eU_nXxeWjODOhimeMyUE47osKXwYtnY4kzkXlUEByoxsl3QExxTq5ZkCRuCZF2A_3fsVR_Y3VwhEdClwx40vLkoVDJNZC6LLck8KLiwrVAgdCUBYTBDW6U3oqJl1m_W1kdiNOT7YU3gMYEgj2qRa0zmL-c3qmKwd7wDuFN4pLx1k0t8XojK5z3f9VKnbAZH5J8FqMLwIf0vP4A79gMFjuSsyWZpWQgFs9-5I0zC5DKYAopme4eJCMHEfc1szT9KLBd1rEDFAoucwxD5bslNM5-Y3XmdMhTHD5uD3hSUZwxFMJCCbPBuEFE3zgzpgils__yhRsEw,,",
        "ndpd-f": "1.w-855182.1.2.Pum0tGVaTU6lJVkg_cEbrg,,.fRB478CKqCbshx2_MERHYE6mvNcUHy5dj_lSb0n_mfEoddz5b1jiX9yc32k2QRogcdAbQZQJHanFgGHZtlh9b7yrf9I3TEXk30rsaLGzFLJFR96xqz7TbnPWRiAr6vzO2dndArgSJ1-2YQ0alfku7V6BvdF1-i7X_kfGX9BPQS6v1Awb2FytNMxgHwIx9p4ziGpxMKNQEPi3-zZy0i1XTKCe06jOqCZKlT5aW8GwoEIZu8H-habZrAqfME123_FysD88Vg_zJC03epKGkBGm7dysBBp4TUFZlVYcbeAKKZulsjLRIMM6x0luKz5-cXIhoHYPHfFvVfKggBGo-XYWBQ,,",
        "2.0.1.1.3.0.7.11.3.1.0.5.11.1.0.5.5": "CN",
        "xAppleActionSignature": "Alf9bFav9DHmnzSlZ8uRkFkd6KVFcX6Kp/3yJGBlCJw3AAACEAQAAAACAAABAKvN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76vN76sAAAAZuRrqGpEsaj9WIO+hiUPO+UVkHIS0/T0JBwAAAKMCpNE6eXa4cgtiKMINZKmYpXuehSUAAACKCQWm018kRh0qwAbGKbnC/6ogVue70AAAANL0L7SoFT5WJgg92oM8qs+Mfex8AQG0qMLptMPWAMtIcq66KeaXhJJzdJv1kJGh6ZJhHWg28Gd1SQVvSO3ODs8i6iZwLGvh9hJo6x2F2GY2LG1xqmI5KBkURs5V9AzjiYPmicaM5u02dLNBsNOz7MwPAAAAMQGu9gDB5A3KAcqV0m46vqDHuqRY0Kfyh1dBP9FwnjGihqvLhoI3ne9DT08UlRaD9FYAAAAAAAAAAAAAAAAAAAEEAgkJAA==",
        "ndpd-bi": "b1.1280x720 1280x660 16 16.-480.en-us",
        "ndpd-vk": "19634",
        "2.0.1.1.3.0.7.11.3.1.0.5.23.5.1.3.1": "宦灵",
        "2.0.1.1.3.0.7.11.3.1.0.5.23.5.7.13.9": "261000",
        "2.0.1.1.3.0.7.11.3.1.0.5.23.5.7.13.7": "潍坊市",
        "state": "9",
        "2.0.1.1.3.0.7.11.3.1.0.5.23.5.7.13.17": "13146231953",
        "2.0.1.1.3.0.7.11.3.1.0.5.23.5.1.3.3": "殊",
        "ndpd-di": "d1-1bf89939d3aa261",
    }

    fmt.Println(elems)
    for _, l := range g.Generate(elems, form).Split(";") {
        fmt.Println(l)
    }
}

var resp = `
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.1//EN" "http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd">
<html xmlns="http://www.apple.com/itms/" lang="zh">


<head>

    <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
    <meta name="keywords" content="iTunes Store" />
    <meta name="description" content="iTunes Store" />
    <meta name="platform-cache-id" content="32" />



    <title>提供付款方式</title>





    <script>
        if (!window.its) {
            window.its = {};
        }
        its.markupLoadStartTime = new Date().getTime();

        if (!window.onerror) window.onerror = function(message, url, lineNumber) {
            if (!window._earlyOnerrorException) {
                window._earlyOnerrorException = {
                    message: message,
                    url: url,
                    lineNumber: lineNumber
                };
            }
        };
    </script>



    <link charset="utf-8" rel="stylesheet" type="text/css" href="https://s.mzstatic.com/htmlResources/8586/desktop-finance-base.css" />
    <link charset="utf-8" rel="stylesheet" type="text/css" href="https://s.mzstatic.com/htmlResources/8586/desktop-finance-itunesstore.css" />
    <link charset="utf-8" rel="stylesheet" type="text/css" href="https://s.mzstatic.com/htmlResources/8586/desktop-finance-modern_ui.css" />
    <link charset="utf-8" rel="stylesheet" type="text/css" href="https://s.mzstatic.com/htmlResources/8586/desktop-finance-billingpage.css" />
    <link charset="utf-8" rel="stylesheet" type="text/css" href="https://s.mzstatic.com/htmlResources/8586/desktop-finance-cup.css" />





    <script type="text/javascript" charset="utf-8" src="/htmlResources/8586/frameworks-primaryinit01.js"></script>





    <script type="text/javascript" charset="utf-8">
        its.serverData = {
            "storePlatformData": {},
            "pageData": {
                "DI.IsPersonalizedDownloadButtonsEnabled": true,
                "DI.PersonalizedDownloadButtonsDKs": "36,37,25,11"
            },
            "constants": {
                "DTBuyButtonMetadataUrl": "https://se.itunes.apple.com/WebObjects/MZStoreElements.woa/wa/buyButtonMetaData?version=2",
                "GlobalLists": {
                    "removeFromWishlistUrl": "https://se.itunes.apple.com/WebObjects/MZStoreElements.woa/wa/removeFromWishlist",
                    "popover": {
                        "width": "350",
                        "maxWishlistItemCount": 20,
                        "height": "70%"
                    },
                    "viewUrl": "https://itunes.apple.com/WebObjects/MZStore.woa/wa/viewGlobalLists?revNum=8586",
                    "updatePreviewHistoryUrl": "https://se.itunes.apple.com/WebObjects/MZStoreElements.woa/wa/updatePreviewHistory",
                    "updateSiriTagHistoryUrl": "https://se.itunes.apple.com/WebObjects/MZStoreElements.woa/wa/updateSiriTagHistory",
                    "viewWishlistUrl": "https://se.itunes.apple.com/WebObjects/MZStoreElements.woa/wa/viewWishlist?cc=cn",
                    "updateRadioHistoryUrl": "https://se.itunes.apple.com/WebObjects/MZStoreElements.woa/wa/updateRadioHistory"
                },
                "Urls": {
                    "upgradeITunesUrl": "http://www.apple.com/itunes/download/",
                    "upgradeSafariUrl": "http://www.apple.com/safari/download/"
                },
                "MZMediaType": {
                    "MobileSoftware": {
                        "id": 8,
                        "name": "MobileSoftwareApplications"
                    },
                    "MacSoftware": {
                        "id": 12,
                        "name": "DesktopApps"
                    },
                    "Audiobooks": {
                        "id": 3,
                        "name": "Audiobooks"
                    },
                    "MusicVideos": {
                        "id": 5,
                        "name": "MusicVideos"
                    },
                    "iTunesU": {
                        "id": 10,
                        "name": "iTunesU"
                    },
                    "EBooks": {
                        "id": 11,
                        "name": "EBooks"
                    },
                    "Music": {
                        "id": 1,
                        "name": "Music"
                    },
                    "Podcasts": {
                        "id": 2,
                        "name": "Podcasts"
                    },
                    "ClassicSoftware": {
                        "id": 7,
                        "name": "iPodGames"
                    },
                    "TVShows": {
                        "id": 4,
                        "name": "TVShows"
                    },
                    "Movies": {
                        "id": 6,
                        "name": "Movies"
                    },
                    "Ringtones": {
                        "id": 9,
                        "name": "Ringtones"
                    },
                    "Textbooks": {
                        "id": 13,
                        "name": "Textbooks"
                    }
                },
                "SFSortOrder": {
                    "PlaylistName": {
                        "id": 9
                    },
                    "RecentBestSellers": {
                        "id": 18
                    },
                    "SeriesOrder": {
                        "id": 17
                    },
                    "CmcRecommended": {
                        "id": 10
                    },
                    "Popularity": {
                        "id": 3
                    },
                    "ArtistName": {
                        "id": 8
                    },
                    "Featured": {
                        "id": 1
                    },
                    "TopRated": {
                        "id": 5
                    },
                    "Duration": {
                        "id": 11
                    },
                    "ReleaseDate": {
                        "id": 2
                    },
                    "TrackNumber": {
                        "id": 13
                    },
                    "Name": {
                        "id": 0
                    },
                    "AllTimeBestSellers": {
                        "id": 14
                    },
                    "ExpirationDate": {
                        "id": 6
                    },
                    "Price": {
                        "id": 7
                    },
                    "MostRecent": {
                        "id": 16
                    },
                    "PurchaseDate": {
                        "id": 4
                    },
                    "DateAdded": {
                        "id": 15
                    },
                    "DiscNumber": {
                        "id": 12
                    }
                },
                "SFEntityType": {
                    "Artist": {
                        "id": 1
                    },
                    "Playlist": {
                        "id": 4
                    },
                    "TopContents": {
                        "id": 7
                    },
                    "MultiRoom": {
                        "id": 3
                    },
                    "Search": {
                        "id": 6
                    },
                    "WeMix": {
                        "id": 8
                    },
                    "Cmc": {
                        "id": 2
                    },
                    "Work": {
                        "id": 10
                    },
                    "Wishlist": {
                        "id": 9
                    },
                    "TopAuthors": {
                        "id": 12
                    },
                    "Room": {
                        "id": 5
                    },
                    "TheyMix": {
                        "id": 11
                    }
                },
                "countryCode": "chn",
                "languageCode": "zh",
                "SFCustomComponentCountryCode": "chn",
                "IXDisplayableKind": {
                    "Artist": {
                        "id": 20,
                        "cssClasses": [],
                        "name": "artist"
                    },
                    "MobileSoftware": {
                        "id": 11,
                        "cssClasses": ["application"],
                        "name": "iosSoftware",
                        "mtId": "8"
                    },
                    "Album": {
                        "id": 2,
                        "cssClasses": ["album", "music"],
                        "name": "album",
                        "mtId": "1"
                    },
                    "ShortFilm": {
                        "id": 8,
                        "cssClasses": ["short-film", "movie", "video"],
                        "name": "shortFilm",
                        "mtId": "6"
                    },
                    "PodcastEpisode": {
                        "id": 15,
                        "cssClasses": ["podcast-episode"],
                        "name": "podcastEpisode",
                        "mtId": "2"
                    },
                    "ApplePubEBook": {
                        "id": 36,
                        "cssClasses": ["ebook"],
                        "name": "ibook",
                        "mtId": "13"
                    },
                    "Song": {
                        "id": 1,
                        "cssClasses": ["song", "music"],
                        "name": "song",
                        "mtId": "1"
                    },
                    "MovieBundle": {
                        "id": 24,
                        "cssClasses": ["bundle", "movie", "video"],
                        "name": "movieBundle",
                        "mtId": "6"
                    },
                    "RingtoneAlbum": {
                        "id": 22,
                        "cssClasses": [],
                        "name": "ringtoneAlbum",
                        "mtId": "9"
                    },
                    "Book": {
                        "id": 3,
                        "cssClasses": ["audiobook"],
                        "name": "book",
                        "mtId": "3"
                    },
                    "ToneAlbum": {
                        "id": 32,
                        "cssClasses": [],
                        "name": "toneAlbum",
                        "mtId": "9"
                    },
                    "TVEpisode": {
                        "id": 6,
                        "cssClasses": ["tv-episode", "tv", "video"],
                        "name": "tvEpisode",
                        "mtId": "4"
                    },
                    "TVSeason": {
                        "id": 7,
                        "cssClasses": ["tv-season", "tv", "video"],
                        "name": "tvSeason",
                        "mtId": "4"
                    },
                    "Course": {
                        "id": 34,
                        "cssClasses": ["course", "itunes-u"],
                        "name": "course",
                        "mtId": "10"
                    },
                    "ClassicSoftwarePackage": {
                        "id": 18,
                        "cssClasses": ["ipod-game"],
                        "name": "gamePackage",
                        "mtId": "7"
                    },
                    "Movie": {
                        "id": 9,
                        "cssClasses": ["movie", "video"],
                        "name": "movie",
                        "mtId": "6"
                    },
                    "Booklet": {
                        "id": 13,
                        "cssClasses": ["booklet", "music"],
                        "name": "booklet",
                        "mtId": "1"
                    },
                    "SocialArtist": {
                        "id": 28,
                        "cssClasses": [],
                        "name": "socialArtist"
                    },
                    "EBook": {
                        "id": 25,
                        "cssClasses": ["ebook"],
                        "name": "epubBook",
                        "mtId": "11"
                    },
                    "MacSoftware": {
                        "id": 30,
                        "cssClasses": ["application", "mac-application"],
                        "name": "desktopApp",
                        "mtId": "12"
                    },
                    "MusicVideo": {
                        "id": 5,
                        "cssClasses": ["music-video", "music", "video"],
                        "name": "musicVideo",
                        "mtId": "5"
                    },
                    "Ringtone": {
                        "id": 21,
                        "cssClasses": [],
                        "name": "ringtone",
                        "mtId": "9"
                    },
                    "Podcast": {
                        "id": 4,
                        "cssClasses": ["podcast"],
                        "name": "podcast",
                        "mtId": "2"
                    },
                    "Concert": {
                        "id": 26,
                        "cssClasses": [],
                        "name": "concert",
                        "mtId": "1"
                    },
                    "ApplePubTextbook": {
                        "id": 37,
                        "cssClasses": ["ebook"],
                        "name": "ibookTextbook",
                        "mtId": "11"
                    },
                    "SoftwareAddOn": {
                        "id": 16,
                        "cssClasses": ["software-add-on", "application"],
                        "name": "softwareAddOn",
                        "mtId": "8"
                    },
                    "Tone": {
                        "id": 31,
                        "cssClasses": [],
                        "name": "tone",
                        "mtId": "9"
                    },
                    "WeMix": {
                        "id": 19,
                        "cssClasses": ["wemix", "music"],
                        "name": "itunesMix",
                        "mtId": "1"
                    },
                    "SocialPerson": {
                        "id": 27,
                        "cssClasses": [],
                        "name": "socialPerson"
                    },
                    "ClassicSoftware": {
                        "id": 10,
                        "cssClasses": ["ipod-game"],
                        "name": "ipodGame",
                        "mtId": "7"
                    },
                    "MetaEBook": {
                        "id": 35,
                        "cssClasses": ["ebook"],
                        "name": "metaEbook",
                        "mtId": "13"
                    },
                    "iTunesPass": {
                        "id": 12,
                        "cssClasses": ["itunes-pass", "album", "music"],
                        "name": "itunesPass",
                        "mtId": "1"
                    },
                    "Mix": {
                        "id": 14,
                        "cssClasses": ["mix", "music"],
                        "name": "mix",
                        "mtId": "1"
                    },
                    "TheyMix": {
                        "id": 33,
                        "cssClasses": [],
                        "name": "thirdPartyMix"
                    }
                },
                "AddToWishlistResult": {
                    "ItemIsFree": {
                        "id": 6
                    },
                    "ListFull": {
                        "id": 3
                    },
                    "SomeContentAlreadyOwnedButUnfulfilled": {
                        "id": 5
                    },
                    "SomeContentAlreadyInWishlist": {
                        "id": 1
                    },
                    "Success": {
                        "id": 0
                    },
                    "Failure": {
                        "id": 2
                    },
                    "SomeContentAlreadyOwned": {
                        "id": 4
                    }
                },
                "getEmailUrl": "https://itunes.apple.com/WebObjects/MZStore.woa/wa/generateEmail?cc=cn",
                "resourceImagePathFrameworks": "https://s.mzstatic.com/htmlResources/8586/frameworks/images/",
                "resourceBasePath": "https://s.mzstatic.com/htmlResources/8586"
            },
            "properties": {
                "clientStatsLoadTimeGroup": "4020",
                "MZHtmlResourcesUtil.allowDeferJsLoad": true,
                "ITSServerEnvironment": "prod",
                "ITSServerInstance": "262462",
                "ITSResourceRevNum": "8586",
                "resourceUrlPrefix": "https://s.mzstatic.com",
                "ITSLogger.ServerReportingProtocol": "https",
                "ITSLogger.ServerReportingDomain": "metrics.mzstatic.com",
                "ITSLogger.ServerReportingApp": "MZUserXP",
                "ITSLogger.RecordStatsAction": "recordStats",
                "ITSLogger.SenderName": "ITSClient",
                "ITSLogger.mirrorToServerByDefaultOnErrors": true,
                "DynaLoader.allowDynaLoading": true,
                "itsLoggerQueueProcessingInterval": 10000,
                "cobaltBundleId": "com.apple.itunesu",
                "getCobaltAppLink": "https://itunes.apple.com/cn/app/id490217893",
                "iTunesUEnrollLink": "https://itunes.apple.com/WebObjects/DZR.woa/wa/iTunesEnrollPopover?cc=cn",
                "cobaltLearnMore": "https://itunes.apple.com/cn/learn-more?about=iTunesUUpgradePage&type=2",
                "isCobaltEnabled": true,
                "isCobaltJavascriptRedirectEnabled": true,
                "isCobaltUpsellPageEnabled": true,
                "personalizedButtonsEnabled": true,
                "SF6.Personalization.isCMAEnabled": false,
                "SF6.Personalization.isCMSEnabled": false,
                "SFJ.isNFAEnabledInStorefront": true,
                "SFJ.isRecommendationsSwooshEnabledInStorefront": true,
                "SF6.StorePlatform.whitelistParams": ["caller", "dsid", "id", "p"],
                "vendLastSupportedVersionIfAvailable": true,
                "vendLastSupportedVersionMacOs": true,
                "isBookGiftingEnabled": true,
                "metrics": {
                    "metricsUrl": "https://xp.apple.com/report",
                    "compoundSeparator": "_",
                    "tokenSeparator": "|",
                    "postFrequency": 60000,
                    "blacklistedFields": ["capacitySystem", "capacitySystemAvailable", "capacityDisk", "capacityData", "capacityDataAvailable"],
                    "impressions": {
                        "viewableThreshold": 1000
                    },
                    "fieldsMap": {
                        "single": {
                            "targetId": ["id", "adamId", "contentId", "type", "fcId", "userPreference", "label", "station-hash", "linkIdentifier"]
                        },
                        "custom": {
                            "location": ["id", "adamId", "dataSetId", "name", "fcKind", "kindIds", "type", "station-hash", "core-seed-name"],
                            "impressions": ["id", "adamId", "station-hash"]
                        },
                        "cookies": ["itcCt", "itscc"],
                        "multi": {}
                    },
                    "metricsBase": {
                        "storeFrontHeader": "143465-19,32",
                        "language": "19",
                        "platformId": "32",
                        "platformName": "iTunes122",
                        "storeFront": "143465",
                        "environmentDataCenter": "NWK"
                    },
                    "Software": {
                        "metricsUrl": "https://xp.apple.com/report",
                        "compoundSeparator": "_",
                        "tokenSeparator": "|",
                        "postFrequency": 60000,
                        "blacklistedFields": ["capacitySystem", "capacitySystemAvailable", "capacityDisk", "capacityData", "capacityDataAvailable"],
                        "impressions": {
                            "viewableThreshold": 1000
                        },
                        "fieldsMap": {
                            "single": {
                                "targetId": ["id", "adamId", "contentId", "type", "fcId", "userPreference", "label", "station-hash", "linkIdentifier"]
                            },
                            "custom": {
                                "location": ["id", "adamId", "dataSetId", "name", "fcKind", "kindIds", "type", "station-hash", "core-seed-name"],
                                "impressions": ["id", "adamId", "station-hash"]
                            },
                            "cookies": ["itcCt", "itscc"],
                            "multi": {}
                        },
                        "metricsBase": {
                            "storeFrontHeader": "143465-19,32",
                            "language": "19",
                            "platformId": "32",
                            "platformName": "iTunes122",
                            "storeFront": "143465",
                            "environmentDataCenter": "NWK"
                        }
                    },
                    "ATV": {
                        "metricsUrl": "https://xp.apple.com/report",
                        "compoundSeparator": "_",
                        "tokenSeparator": "|",
                        "postFrequency": 60000,
                        "blacklistedFields": ["capacitySystem", "capacitySystemAvailable", "capacityDisk", "capacityData", "capacityDataAvailable"],
                        "impressions": {
                            "viewableThreshold": 1000
                        },
                        "fieldsMap": {
                            "single": {
                                "targetId": ["id", "adamId", "contentId", "type", "fcId", "userPreference", "label", "station-hash", "linkIdentifier"]
                            },
                            "custom": {
                                "location": ["id", "adamId", "dataSetId", "name", "fcKind", "kindIds", "type", "station-hash", "core-seed-name"],
                                "impressions": ["id", "adamId", "station-hash"]
                            },
                            "cookies": ["itcCt", "itscc"],
                            "multi": {}
                        },
                        "metricsBase": {
                            "storeFrontHeader": "143465-19,32",
                            "language": "19",
                            "platformId": "32",
                            "platformName": "iTunes122",
                            "storeFront": "143465",
                            "environmentDataCenter": "NWK"
                        }
                    },
                    "eBooks": {
                        "metricsUrl": "https://xp.apple.com/report",
                        "compoundSeparator": "_",
                        "tokenSeparator": "|",
                        "postFrequency": 60000,
                        "blacklistedFields": ["capacitySystem", "capacitySystemAvailable", "capacityDisk", "capacityData", "capacityDataAvailable"],
                        "impressions": {
                            "viewableThreshold": 1000
                        },
                        "fieldsMap": {
                            "single": {
                                "targetId": ["id", "adamId", "contentId", "type", "fcId", "userPreference", "label", "station-hash", "linkIdentifier"]
                            },
                            "custom": {
                                "location": ["id", "adamId", "dataSetId", "name", "fcKind", "kindIds", "type", "station-hash", "core-seed-name"],
                                "impressions": ["id", "adamId", "station-hash"]
                            },
                            "cookies": ["itcCt", "itscc"],
                            "multi": {}
                        },
                        "metricsBase": {
                            "storeFrontHeader": "143465-19,32",
                            "language": "19",
                            "platformId": "32",
                            "platformName": "iTunes122",
                            "storeFront": "143465",
                            "environmentDataCenter": "NWK"
                        }
                    },
                    "iTunes": {
                        "metricsUrl": "https://xp.apple.com/report",
                        "compoundSeparator": "_",
                        "tokenSeparator": "|",
                        "postFrequency": 60000,
                        "blacklistedFields": ["capacitySystem", "capacitySystemAvailable", "capacityDisk", "capacityData", "capacityDataAvailable"],
                        "impressions": {
                            "viewableThreshold": 1000
                        },
                        "fieldsMap": {
                            "single": {
                                "targetId": ["id", "adamId", "contentId", "type", "fcId", "userPreference", "label", "station-hash", "linkIdentifier"]
                            },
                            "custom": {
                                "location": ["id", "adamId", "dataSetId", "name", "fcKind", "kindIds", "type", "station-hash", "core-seed-name"],
                                "impressions": ["id", "adamId", "station-hash"]
                            },
                            "cookies": ["itcCt", "itscc"],
                            "multi": {}
                        },
                        "metricsBase": {
                            "storeFrontHeader": "143465-19,32",
                            "language": "19",
                            "platformId": "32",
                            "platformName": "iTunes122",
                            "storeFront": "143465",
                            "environmentDataCenter": "NWK"
                        }
                    },
                    "iTunesU": {
                        "metricsUrl": "https://xp.apple.com/report",
                        "compoundSeparator": "_",
                        "tokenSeparator": "|",
                        "postFrequency": 60000,
                        "blacklistedFields": ["capacitySystem", "capacitySystemAvailable", "capacityDisk", "capacityData", "capacityDataAvailable"],
                        "impressions": {
                            "viewableThreshold": 1000
                        },
                        "fieldsMap": {
                            "single": {
                                "targetId": ["id", "adamId", "contentId", "type", "fcId", "userPreference", "label", "station-hash", "linkIdentifier"]
                            },
                            "custom": {
                                "location": ["id", "adamId", "dataSetId", "name", "fcKind", "kindIds", "type", "station-hash", "core-seed-name"],
                                "impressions": ["id", "adamId", "station-hash"]
                            },
                            "cookies": ["itcCt", "itscc"],
                            "multi": {}
                        },
                        "metricsBase": {
                            "storeFrontHeader": "143465-19,32",
                            "language": "19",
                            "platformId": "32",
                            "platformName": "iTunes122",
                            "storeFront": "143465",
                            "environmentDataCenter": "NWK"
                        }
                    },
                    "MacAppStore": {
                        "metricsUrl": "https://xp.apple.com/report",
                        "compoundSeparator": "_",
                        "tokenSeparator": "|",
                        "postFrequency": 60000,
                        "blacklistedFields": ["capacitySystem", "capacitySystemAvailable", "capacityDisk", "capacityData", "capacityDataAvailable"],
                        "impressions": {
                            "viewableThreshold": 1000
                        },
                        "fieldsMap": {
                            "single": {
                                "targetId": ["id", "adamId", "contentId", "type", "fcId", "userPreference", "label", "station-hash", "linkIdentifier"]
                            },
                            "custom": {
                                "location": ["id", "adamId", "dataSetId", "name", "fcKind", "kindIds", "type", "station-hash", "core-seed-name"],
                                "impressions": ["id", "adamId", "station-hash"]
                            },
                            "cookies": ["itcCt", "itscc"],
                            "multi": {}
                        },
                        "metricsBase": {
                            "storeFrontHeader": "143465-19,32",
                            "language": "19",
                            "platformId": "32",
                            "platformName": "iTunes122",
                            "storeFront": "143465",
                            "environmentDataCenter": "NWK"
                        }
                    },
                    "Maps": {
                        "metricsUrl": "https://xp.apple.com/report",
                        "compoundSeparator": "_",
                        "tokenSeparator": "|",
                        "postFrequency": 60000,
                        "blacklistedFields": ["capacitySystem", "capacitySystemAvailable", "capacityDisk", "capacityData", "capacityDataAvailable"],
                        "impressions": {
                            "viewableThreshold": 1000
                        },
                        "fieldsMap": {
                            "single": {
                                "targetId": ["id", "adamId", "contentId", "type", "fcId", "userPreference", "label", "station-hash", "linkIdentifier"]
                            },
                            "custom": {
                                "location": ["id", "adamId", "dataSetId", "name", "fcKind", "kindIds", "type", "station-hash", "core-seed-name"],
                                "impressions": ["id", "adamId", "station-hash"]
                            },
                            "cookies": ["itcCt", "itscc"],
                            "multi": {}
                        },
                        "metricsBase": {
                            "storeFrontHeader": "143465-19,32",
                            "language": "19",
                            "platformId": "32",
                            "platformName": "iTunes122",
                            "storeFront": "143465",
                            "environmentDataCenter": "NWK"
                        }
                    },
                    "MusicPlayer": {
                        "metricsUrl": "https://xp.apple.com/report",
                        "compoundSeparator": "_",
                        "tokenSeparator": "|",
                        "postFrequency": 60000,
                        "blacklistedFields": ["capacitySystem", "capacitySystemAvailable", "capacityDisk", "capacityData", "capacityDataAvailable"],
                        "impressions": {
                            "viewableThreshold": 1000
                        },
                        "fieldsMap": {
                            "single": {
                                "targetId": ["id", "adamId", "contentId", "type", "fcId", "userPreference", "label", "station-hash", "linkIdentifier"]
                            },
                            "custom": {
                                "location": ["id", "adamId", "dataSetId", "name", "fcKind", "kindIds", "type", "station-hash", "core-seed-name"],
                                "impressions": ["id", "adamId", "station-hash"]
                            },
                            "cookies": ["itcCt", "itscc"],
                            "multi": {}
                        },
                        "metricsBase": {
                            "storeFrontHeader": "143465-19,32",
                            "language": "19",
                            "platformId": "32",
                            "platformName": "iTunes122",
                            "storeFront": "143465",
                            "environmentDataCenter": "NWK"
                        }
                    },
                    "Podcasts": {
                        "metricsUrl": "https://xp.apple.com/report",
                        "compoundSeparator": "_",
                        "tokenSeparator": "|",
                        "postFrequency": 60000,
                        "blacklistedFields": ["capacitySystem", "capacitySystemAvailable", "capacityDisk", "capacityData", "capacityDataAvailable"],
                        "impressions": {
                            "viewableThreshold": 1000
                        },
                        "fieldsMap": {
                            "single": {
                                "targetId": ["id", "adamId", "contentId", "type", "fcId", "userPreference", "label", "station-hash", "linkIdentifier"]
                            },
                            "custom": {
                                "location": ["id", "adamId", "dataSetId", "name", "fcKind", "kindIds", "type", "station-hash", "core-seed-name"],
                                "impressions": ["id", "adamId", "station-hash"]
                            },
                            "cookies": ["itcCt", "itscc"],
                            "multi": {}
                        },
                        "metricsBase": {
                            "storeFrontHeader": "143465-19,32",
                            "language": "19",
                            "platformId": "32",
                            "platformName": "iTunes122",
                            "storeFront": "143465",
                            "environmentDataCenter": "NWK"
                        }
                    },
                    "Third-Party": {
                        "metricsUrl": "https://xp.apple.com/report",
                        "compoundSeparator": "_",
                        "tokenSeparator": "|",
                        "postFrequency": 60000,
                        "blacklistedFields": ["capacitySystem", "capacitySystemAvailable", "capacityDisk", "capacityData", "capacityDataAvailable"],
                        "impressions": {
                            "viewableThreshold": 1000
                        },
                        "fieldsMap": {
                            "single": {
                                "targetId": ["id", "adamId", "contentId", "type", "fcId", "userPreference", "label", "station-hash", "linkIdentifier"]
                            },
                            "custom": {
                                "location": ["id", "adamId", "dataSetId", "name", "fcKind", "kindIds", "type", "station-hash", "core-seed-name"],
                                "impressions": ["id", "adamId", "station-hash"]
                            },
                            "cookies": ["itcCt", "itscc"],
                            "multi": {}
                        },
                        "metricsBase": {
                            "storeFrontHeader": "143465-19,32",
                            "language": "19",
                            "platformId": "32",
                            "platformName": "iTunes122",
                            "storeFront": "143465",
                            "environmentDataCenter": "NWK"
                        }
                    },
                    "VideoPlayer": {
                        "metricsUrl": "https://xp.apple.com/report",
                        "compoundSeparator": "_",
                        "tokenSeparator": "|",
                        "postFrequency": 60000,
                        "blacklistedFields": ["capacitySystem", "capacitySystemAvailable", "capacityDisk", "capacityData", "capacityDataAvailable"],
                        "impressions": {
                            "viewableThreshold": 1000
                        },
                        "fieldsMap": {
                            "single": {
                                "targetId": ["id", "adamId", "contentId", "type", "fcId", "userPreference", "label", "station-hash", "linkIdentifier"]
                            },
                            "custom": {
                                "location": ["id", "adamId", "dataSetId", "name", "fcKind", "kindIds", "type", "station-hash", "core-seed-name"],
                                "impressions": ["id", "adamId", "station-hash"]
                            },
                            "cookies": ["itcCt", "itscc"],
                            "multi": {}
                        },
                        "metricsBase": {
                            "storeFrontHeader": "143465-19,32",
                            "language": "19",
                            "platformId": "32",
                            "platformName": "iTunes122",
                            "storeFront": "143465",
                            "environmentDataCenter": "NWK"
                        }
                    },
                    "WiFi-Music": {
                        "metricsUrl": "https://xp.apple.com/report",
                        "compoundSeparator": "_",
                        "tokenSeparator": "|",
                        "postFrequency": 60000,
                        "blacklistedFields": ["capacitySystem", "capacitySystemAvailable", "capacityDisk", "capacityData", "capacityDataAvailable"],
                        "impressions": {
                            "viewableThreshold": 1000
                        },
                        "fieldsMap": {
                            "single": {
                                "targetId": ["id", "adamId", "contentId", "type", "fcId", "userPreference", "label", "station-hash", "linkIdentifier"]
                            },
                            "custom": {
                                "location": ["id", "adamId", "dataSetId", "name", "fcKind", "kindIds", "type", "station-hash", "core-seed-name"],
                                "impressions": ["id", "adamId", "station-hash"]
                            },
                            "cookies": ["itcCt", "itscc"],
                            "multi": {}
                        },
                        "metricsBase": {
                            "storeFrontHeader": "143465-19,32",
                            "language": "19",
                            "platformId": "32",
                            "platformName": "iTunes122",
                            "storeFront": "143465",
                            "environmentDataCenter": "NWK"
                        }
                    },
                    "AppStore": {
                        "metricsUrl": "https://xp.apple.com/report",
                        "compoundSeparator": "_",
                        "tokenSeparator": "|",
                        "postFrequency": 60000,
                        "blacklistedFields": ["capacitySystem", "capacitySystemAvailable", "capacityDisk", "capacityData", "capacityDataAvailable"],
                        "impressions": {
                            "viewableThreshold": 1000
                        },
                        "fieldsMap": {
                            "single": {
                                "targetId": ["id", "adamId", "contentId", "type", "fcId", "userPreference", "label", "station-hash", "linkIdentifier"]
                            },
                            "custom": {
                                "location": ["id", "adamId", "dataSetId", "name", "fcKind", "kindIds", "type", "station-hash", "core-seed-name"],
                                "impressions": ["id", "adamId", "station-hash"]
                            },
                            "cookies": ["itcCt", "itscc"],
                            "multi": {}
                        },
                        "metricsBase": {
                            "storeFrontHeader": "143465-19,32",
                            "language": "19",
                            "platformId": "32",
                            "platformName": "iTunes122",
                            "storeFront": "143465",
                            "environmentDataCenter": "NWK"
                        }
                    },
                    "MobileStore": {
                        "metricsUrl": "https://xp.apple.com/report",
                        "compoundSeparator": "_",
                        "tokenSeparator": "|",
                        "postFrequency": 60000,
                        "blacklistedFields": ["capacitySystem", "capacitySystemAvailable", "capacityDisk", "capacityData", "capacityDataAvailable"],
                        "impressions": {
                            "viewableThreshold": 1000
                        },
                        "fieldsMap": {
                            "single": {
                                "targetId": ["id", "adamId", "contentId", "type", "fcId", "userPreference", "label", "station-hash", "linkIdentifier"]
                            },
                            "custom": {
                                "location": ["id", "adamId", "dataSetId", "name", "fcKind", "kindIds", "type", "station-hash", "core-seed-name"],
                                "impressions": ["id", "adamId", "station-hash"]
                            },
                            "cookies": ["itcCt", "itscc"],
                            "multi": {}
                        },
                        "metricsBase": {
                            "storeFrontHeader": "143465-19,32",
                            "language": "19",
                            "platformId": "32",
                            "platformName": "iTunes122",
                            "storeFront": "143465",
                            "environmentDataCenter": "NWK"
                        }
                    }
                },
                "MZStore.shouldAttemptToEvaluateXpathSecondTime": true,
                "MZStore.shouldCacheElementTypeProtoObjects": true,
                "isStoweUpgradeDialogEnabled": false,
                "stoweUpgradeSize": 1600000000
            }
        }
    </script>
    <script type="text/javascript" charset="utf-8">
        if (!window.iTSLocalization) {
            iTSLocalization = {};
        }
        iTSLocalization._strings = {
            "Js.MyAlerts.GeniusActivationNeeded.Button": "进入 Genius",
            "Js.Wishlist.RemoveTooltip": "从欲购清单上删除",
            "Js.DI.HDCP.UpgradeiTunes.MovieBundle.explanation": "要购买此项目，请升级至 iTunes 的最新版本。",
            "Js.Showcase.Next": "下一个",
            "Js.DI.HDCP.UpgradeiTunes.Movie.explanation": "要购买此项目，请升级至 iTunes 的最新版本。",
            "Js.CNConnections.Dialog.GetStarted": "Get Started",
            "Js.ReportAConcern.SigninMessage": "请输入您的 Apple ID 和密码，然后点击“登录”。",
            "DIBuyButtonDownloadingAria": "正在下载, @@name@@",
            "JS.MZFinance.Captcha.ServerTimeout": "无法载入。 请稍后再试。",
            "MZStore.subscribed": "已订阅",
            "Js.TextTruncation.More": "更多",
            "Js.QuickView.Unavailable.Title": "本产品目前无法快速查看。",
            "DIHPurchasesPage.UpgradeiTunes.message": "需要升级",
            "Js.DI.BuyButton.Downloaded": "已下载",
            "JS.MZCommerce.CarrierBilling.TooManyOTPResends": "You’ve requested a verification code too many times and must wait to request a new one. Try again later.",
            "CNConnections.Dialog.Post.Review.SignIn.Title": "Sign in to post.",
            "CNConnections.Upload.UploadErrorHeading": "上传出现问题。",
            "CNConnections.Dialog.StopFollowing.Message": "CNConnections.Dialog.StopFollowing.Message",
            "CNConnections.Dialog.OptIn.Title": "Ping 将从 9 月 30 日起停止使用",
            "Js.DI.CreateAppleId": "创建 Apple ID",
            "DIBuyButtonDownloadedAria": "已下载，@@name@@",
            "Js.MyAlerts.GeniusActivationNeeded.Explanation": "只有当您打开Genius 时，才会启用基于您的资料库内容的提示。语了解更多，请单击“转到 Genius”。",
            "JS.MZCommerce.CarrierBilling.MobileNumberRequired": "需要提供您帐户的电话号码。",
            "_thousandsSeparator": ",",
            "JS.MZCommerce.CarrierBilling.AjaxFailMessage": "Your request could not be completed. Please try again.",
            "InAppPurchases": "App 内购买",
            "JS.MZCommerce.CarrierBilling.IOException": "Your request could not be completed. Please try again later.",
            "Js.CNConnections.LoginRequired.Message": "此操作需要使用 Ping。Ping 将从 9 月 30 日起停止使用，我们不再接受新的 Ping 成员。如果您已经是 Ping 成员，请在下方登录以继续。",
            "Js.Pagination.PreviousPage": "前一页",
            "DTMediaPicker.RemoveSong": "移除歌曲",
            "Js.TellAFriend.SigninTitle": "登录以告诉朋友。",
            "Js.InlineReview.RateThis": "我来评分",
            "Js.InlineReview.5": "好极了",
            "Js.InlineReview.4": "很不错",
            "Js.InlineReview.3": "还可以",
            "Js.InlineReview.2": "不喜欢",
            "JS.MZFinance.Captcha.ReplayAudioCaptcha": "重播验证码",
            "Js.InlineReview.1": "痛恨",
            "DownloadNow": "下载",
            "Js.InlinePreview.StopPreview": "停止播放《@@title@@》预览",
            "DTMediaPicker.AddSongs": "添加歌曲",
            "JS.MZCommerce.CarrierBilling.CarrierNotEligible": "Mobile phone billing is only available with @@launchedCarriersVisibleToPublic@@. Enter a different mobile number or select another payment type.",
            "DIDLargeEbookLockup.AfterDownloadCopy": "本书的试读本已传送到您设备上的、已启用了自动下载功能的 iBooks 中。",
            "JS.MZCommerce.CarrierBilling.RejectedByCarrier": "We are unable to process your request due to an error with your mobile account. Please contact your carrier.",
            "CNConnections.Dialog.OptIn.Message": "CNConnections.Dialog.OptIn.Message",
            "Js.DI.BuyButton.Update": "更新",
            "DTHEditiMixPage.Error.PlaylistMaxReachedRemoveToAdd": "此歌单已达到最多 @@count@@ 首歌曲的上限。您必须先删除一些歌曲，然后才能添加其它歌曲。",
            "Js.CreateAccount": "创建 Apple ID",
            "DIBuyButtonAria": "@@buyActionText@@，@@name@@： @@price@@",
            "Js.CNConnections.OptInRequired.Title": "Ping",
            "Cancel": "取消",
            "Js.ManageArtistAlerts.CheckAll.Actors": "查看所有演员",
            "DIHPurchasesPage.UpgradeSafari.explanation": "必须升级至 Safari 的最新版本才能使用此功能。",
            "More": "更多",
            "Js.DI.BuyButton.Preordered": "已预订",
            "Js.Pagination.DisabledButtonText": "@@button_text@@ 已禁用",
            "Js.MyAlerts.CancelEmailSignUp.Message": "您确定要停止接受电子邮件通报吗？",
            "JS.errors.requiredLite": "您没有完整填写表格。",
            "CNConnections.Dialog.Post.Review.SignIn.Message": "CNConnections.Dialog.Post.Review.SignIn.Message",
            "Js.ManageArtistAlerts.UncheckAll.Artists": "清除勾选所有艺人",
            "Js.DI.BuyButton.Downloading": "正在下载",
            "Js.InlineRating.SigninTitle": "登录对此产品评分。",
            "Js.ManageArtistAlerts.UncheckAll.Actors": "清除勾选所有演员",
            "JS.MZCommerce.CarrierBilling.TooManyOTPRequests": "You have exceeded the number of verification codes you can request at this time. Please try again later.",
            "JS.MZFinance.EditAccount.CaptchaTitle.Audio": "输入您听到的字符：",
            "Js.CNConnections.OptInRequired.Message": "This requires Ping. Ping will no longer be available as of September 30, and we are not accepting new members. If you are already a Ping member, sign in below to continue.",
            "JS.MZCommerce.CarrierBilling.UnknownError": "Your request could not be completed.",
            "JS.MZCommerce.CarrierBilling.InvalidOTP": "The verification code you have entered is invalid. Please check the code and try again.",
            "Js.DI.HDCP.UpgradeiTunes.MovieBundle.message": "无法在该电脑上购买此项目的高清版。",
            "Js.Search.HintsTitle": "建议",
            "JS.MZCommerce.CarrierBilling.InvalidPhoneNumber": "Mobile phone billing is only available with @@launchedCarriersVisibleToPublic@@. Enter a different mobile number or select another payment type.",
            "Js.InlinePreview.PlayPreview": "播放《@@title@@》预览",
            "SignIn": "登录",
            "Js.ManageArtistAlerts.UncheckAll": "清除全部勾选",
            "Js.iTunesStoreError.Message": "无法完成您的请求。",
            "DTHEditiMixPage.Error.PlaylistMaxReached": "此歌单已达到最多 @@count@@ 首歌曲的上限。",
            "Js.Pagination.NextPage": "下页",
            "Js.CNConnections.RequestSent": "Request sent",
            "JS.MZCommerce.CarrierBilling.InternalError": "An error occurred while attempting to process your request. Please try again later.",
            "DIHPurchasesPage.UpgradeSafari.message": "需要升级",
            "JS.MZCommerce.CarrierBilling.IntegratorError": "We are unable to process your request at this time. Please try again later.",
            "CNConnections.Upload.UploadErrorMessage": "抱歉，上传出错。",
            "Js.InlineReview.SigninMessage": "请输入您的 Apple ID 和密码，然后点击“登录”。",
            "DIBuyButtonPreorderedAria": "已预订, @@name@@",
            "Less": "更少",
            "Js.MyAlerts.GeniusActivationNeeded.Message": "这需要打开 Genius。",
            "Js.ManageArtistAlerts.CheckAll.Artists": "查看所有艺人",
            "Js.TellAFriend.SigninMessage": "输入您的 Apple ID 和密码，然后点击“登录”。",
            "JS.MZCommerce.CarrierBilling.CarrierError": "An error occurred while attempting to process your request. Please try again later.",
            "CNConnections.Dialog.StopFollowing.Title": "Are you sure you want to stop following @@fullName@@?",
            "Js.InlineReview.SigninTitle": "登录以撰写评价。",
            "Js.InlineRating.SigninMessage": "请登录以继续操作。",
            "JS.MZCommerce.CarrierBilling.SilentSignupFailed": "We are not able to complete your request. Please try again.",
            "Js.Pagination.Back": "上一页",
            "CNConnections.Dialog.PrivateUser.IllegalAction.Message": "If you would like to perform this action, please change your privacy settings.",
            "SFiTunesUSubscribeDialog_explanation": "您的课程将下载。音频、视频、图书等教材可从您课程笔记本的教材列表中下载或购买。",
            "Js.QuickView.Unavailable.Text": "请稍后再试。",
            "Js.MyAlerts.ConfirmEmailSignUp.Explanation": "您可以随时从您的”我的通报“页面更改这项偏好设置。",
            "Js.CNConnections.Dialog.Cancel": "取消",
            "Js.CNConnections.PeoplePopupMore": "和其他 @@count@@ 位...",
            "Js.CNConnections.Confirmed": "Confirmed",
            "Js.CNConnections.UserReview.YouLiked": "You Liked",
            "JS.MZCommerce.CarrierBilling.InvalidRequest": "Your request could not be completed.",
            "JS.MZFinance.Captcha.Loading": "正在载入字符。",
            "Js.Pagination.Next": "下一页",
            "Js.ManageArtistAlerts.CheckAll": "显示全部",
            "JS.MZApplication.GeneralError_explanation": "请稍后再试。",
            "Js.ReportAConcern.SigninTitle": "请登录以报告问题。",
            "Js.MyAlerts.CancelEmailSignUp.Explanation": "您可以随时从您的”我的通报“页面更改这项偏好设置。",
            "CNConnections.Dialog.PrivateUser.LeakyAction.Message": "This action will make your profile picture and name visible to others.",
            "Js.TextLimit.Remaining": "剩余 @@count@@ 个字符",
            "CNConnections.Dialog.PrivateUser.LeakyAction.Title": "Are you sure you want to perform this action?",
            "Js.Pagination.PageNumberTitle": "第 @@num@@ 页",
            "JS.MZFinance.Captcha.LoadError": "无法载入。 请稍后再试。",
            "JS.MZFinance.EditAccount.CaptchaTitle.Text": "输入动态字符：",
            "JS.MZFinance.Required": "必填",
            "JS.errors.noTitle": "必须选择一个称呼。",
            "MZStore.sampleSentButton": "样书已发送",
            "DIHPurchasesPage.UpgradeiTunes.explanation": "必须升级至 iTunes 的最新版本才能使用此功能。",
            "JS.MZCommerce.CarrierBilling.ExpiredOTP": "The verification code you have entered has expired. Please request a new verification code.",
            "SFiTunesUSubscribeDialog_message": "iTunes U 课程可能会使用需付费的教材。",
            "Js.InlineReview.Thanks": "多谢！",
            "Js.InlineReview.Error": "错误",
            "CNConnections.Dialog.Like.Review.SignIn.Message": "CNConnections.Dialog.Like.Review.SignIn.Message",
            "Js.List.Item": "产品",
            "Js.DI.HDCP.UpgradeiTunes.Movie.message": "无法在该电脑上购买此项目的高清版。",
            "Js.iTunesStoreError.Explanation": "iTunes Store 发生故障，请稍后再试。 (@@errorNum@@)",
            "Js.InlineReview.ClickToRate": "点击评分",
            "CNConnections.Dialog.PrivateUser.IllegalAction.Title": "Private users may not perform this action.",
            "Js.MyAlerts.ConfirmEmailSignUp.Message": "您确定想要收到电子邮件通报吗？",
            "CNConnections.Dialog.OptIn.ButtonSubmit": "了解详情",
            "MHJavascriptLocalizations.DIBuyButtonPurchased": "已购项目",
            "DIBuyButtonPurchasedAria": "已购项目, @@name@@",
            "_decimalSeparator": ".",
            "CNConnections.Dialog.OptIn.ButtonCancel": "取消",
            "DIBuyButtonRedownloadAria": "Download, @@name@@: Free",
            "Js.CNConnections.LoginRequired.Title": "登录以访问 Ping",
            "CNConnections.Dialog.OptIn.ButtonSignIn": "Sign In",
            "JS.MZCommerce.CarrierBilling.TooManyOTPValidates": "You have exceeded the number of code verifications you can have at this time. Please try again later.",
            "Js.Pagination.PageNumber": "第 @@num@@ 页",
            "Js.DI.HDCP.UpgradeiTunes.button": "升级 iTunes",
            "DIBuyButtonUpdateAria": "更新，@@name@@",
            "JS.MZCommerce.CarrierBilling.CarrierServiceDown": "We are currently unable to contact your carrier in order to process your request. Please try again later.",
            "CNConnections.Dialog.Like.Review.SignIn.Title": "Sign in to like a review.",
            "JS.MZCommerce.CarrierBilling.ChooseBillingPhoneOption": "请选择一个帐单电话号码选项。",
            "CNConnections.Upload.MaxPhotosErrorMessage": "一次最多可以上传 20 张照片。",
            "Js.DI.BuyButton.Download": "下载",
            "Js.DI.BuyButton.DownloadAll": "全部下载"
        }
    </script>

    <script type="text/javascript" charset="utf-8" src="/htmlResources/8586/desktop-finance-base.js"></script>
    <script type="text/javascript" charset="utf-8" src="/htmlResources/8586/desktop-finance-itunesstore.js"></script>
    <script type="text/javascript" charset="utf-8" src="/htmlResources/8586/desktop-finance-billingpage.js"></script>
    <script type="text/javascript" charset="utf-8" src="/htmlResources/8586/desktop-finance-cup.js"></script>
    <script type="text/javascript" charset="utf-8" src="/htmlResources/8586/desktop-finance-accountsignature.js"></script>





    <script id="protocol" type="text/x-apple-plist">
        <dict>



            <key>pings</key>
            <array>


            </array>





        </dict>
    </script>



</head>

<body cn-cookie-name="mz_at1" data-country="chn" cn-if-cookie-name="mz_if" disable-history-and-navigation="true" class="chn zh account its edit-address signup no-country-nav simple-footer">
    <script type="text/javascript">
        if (window['iTSKit']) {
            iTSKit.preLoad();
        }
        if (window.its && its.kit) {
            its.kit.preload();
        }
    </script>





    <div class="content">





        <span class="secure-connection">安全连接</span>




        <form name="f_2_0_1_1_3_0_7_11_3" method="post" action="/WebObjects/MZFinance.woa/wo/5.2.0.1.1.3.0.7.11.3">


            <h1>提供付款方式</h1>
            <section class="edit-billing">





                <div class="billing-country-information">





                    如果您现在提供付款方式，则在您进行购买时才会收取费用。如果您未选择付款方式，则在您第一次购买时将要求您提供付款方式。





                </div>


                <!--

    -->


                <!--

    -->

                <div class="change-country">
                    <p>如果付款信息中的帐单地址不在中国，<a href="#" class="show-bubble">点击此处</a>。</p>
                    <div class="bubble">
                        <div class="bubble-wrapper">


                            <p>您需提供在有 iTunes 服务的国家或地区中有效的付款方式和帐单地址。请从下面的列表中选择一个国家/地区，如果某个国家/地区未列出，说明 iTunes 在此地尚未提供。</p>
                            <label for="country">选择一个国家/地区:</label>
                            <div class="select">
                                <select id="country" name="2.0.1.1.3.0.7.11.3.1.0.5.11.1.0.5.5">
<option value="AL">阿尔巴尼亚</option>
<option value="DZ">阿尔及利亚</option>
<option value="AR">阿根廷</option>
<option value="AE">阿联酋</option>
<option value="OM">阿曼</option>
<option value="AZ">阿塞拜疆</option>
<option value="EG">埃及</option>
<option value="IE">爱尔兰</option>
<option value="EE">爱沙尼亚</option>
<option value="AO">安哥拉</option>
<option value="AI">安圭拉岛</option>
<option value="AG">安提瓜和巴布达</option>
<option value="AT">奥地利</option>
<option value="AU">澳大利亚</option>
<option value="BB">巴巴多斯</option>
<option value="PG">巴布亚新几内亚</option>
<option value="BS">巴哈马</option>
<option value="PK">巴基斯坦</option>
<option value="PY">巴拉圭</option>
<option value="BH">巴林</option>
<option value="PA">巴拿马</option>
<option value="BR">巴西</option>
<option value="BY">白俄羅斯</option>
<option value="BM">百慕大</option>
<option value="BG">保加利亚</option>
<option value="BJ">贝宁</option>
<option value="BE">比利时</option>
<option value="IS">冰岛</option>
<option value="BO">玻利维亚</option>
<option value="PL">波兰</option>
<option value="BW">博茨瓦纳</option>
<option value="BZ">伯利兹</option>
<option value="BT">不丹</option>
<option value="BF">布基纳法索</option>
<option value="DK">丹麦</option>
<option value="DE">德国</option>
<option value="DM">多米尼加</option>
<option value="DO">多米尼加</option>
<option value="RU">俄罗斯</option>
<option value="EC">厄瓜多尔</option>
<option value="FR">法国</option>
<option value="PH">菲律宾</option>
<option value="FI">芬兰</option>
<option value="CV">佛得角</option>
<option value="GM">冈比亚</option>
<option value="CG">刚果共和国</option>
<option value="CO">哥伦比亚</option>
<option value="CR">哥斯达黎加</option>
<option value="GD">格林纳达</option>
<option value="GY">圭亚那</option>
<option value="KZ">哈萨克斯坦</option>
<option value="KR">韩国</option>
<option value="NL">荷兰</option>
<option value="HN">洪都拉斯</option>
<option value="KG">吉尔吉兹斯坦</option>
<option value="GW">几内亚比绍</option>
<option value="CA">加拿大</option>
<option value="GH">加纳</option>
<option value="KH">柬埔寨</option>
<option value="CZ">捷克共和国</option>
<option value="ZW">津巴布韦</option>
<option value="QA">卡塔尔</option>
<option value="KY">开曼群岛</option>
<option value="KW">科威特</option>
<option value="HR">克罗地亚</option>
<option value="KE">肯尼亚</option>
<option value="LV">拉脱维亚</option>
<option value="LA">老挝</option>
<option value="LB">黎巴嫩</option>
<option value="LR">利比里亚</option>
<option value="LT">立陶宛</option>
<option value="LU">卢森堡</option>
<option value="RO">罗马尼亚</option>
<option value="MG">马达加斯加</option>
<option value="MT">马耳他</option>
<option value="MW">马拉维</option>
<option value="MY">马来西亚</option>
<option value="ML">马里</option>
<option value="MK">马其顿</option>
<option value="MU">毛里求斯</option>
<option value="MR">毛里塔尼亚</option>
<option value="US">美国</option>
<option value="MN">蒙古</option>
<option value="MS">蒙特塞拉特</option>
<option value="PE">秘鲁</option>
<option value="FM">密克罗尼西亚</option>
<option value="MD">摩尔多瓦</option>
<option value="MZ">莫桑比克</option>
<option value="MX">墨西哥</option>
<option value="NA">纳米比亚</option>
<option value="ZA">南非</option>
<option value="NP">尼泊尔</option>
<option value="NI">尼加拉瓜</option>
<option value="NE">尼日尔州</option>
<option value="NG">尼日利亚</option>
<option value="NO">挪威</option>
<option value="PW">帕劳</option>
<option value="PT">葡萄牙</option>
<option value="JP">日本</option>
<option value="SE">瑞典</option>
<option value="CH">瑞士</option>
<option value="SV">萨尔瓦多</option>
<option value="SL">塞拉利昂</option>
<option value="SN">塞内加尔</option>
<option value="CY">塞浦路斯</option>
<option value="SC">塞舌尔</option>
<option value="SA">沙特阿拉伯</option>
<option value="ST">圣多美和普林西比</option>
<option value="KN">圣基茨岛和尼维斯</option>
<option value="LC">圣卢西亚</option>
<option value="VC">圣文森特和格林纳丁斯</option>
<option value="LK">斯里兰卡</option>
<option value="SK">斯洛伐克</option>
<option value="SI">斯洛文尼亚</option>
<option value="SZ">斯威士兰</option>
<option value="SR">苏里南</option>
<option value="SB">所罗门群岛</option>
<option value="TJ">塔吉克斯坦</option>
<option value="TW">台湾</option>
<option value="TH">泰国</option>
<option value="TZ">坦桑尼亚</option>
<option value="TC">特克斯和凯科斯群岛</option>
<option value="TT">特立尼达和多巴哥</option>
<option value="TN">突尼斯</option>
<option value="TR">土耳其</option>
<option value="TM">土库曼斯坦</option>
<option value="GT">危地马拉</option>
<option value="VG">维尔京群岛（英属）</option>
<option value="VE">委内瑞拉</option>
<option value="BN">文莱</option>
<option value="UG">乌干达</option>
<option value="UA">乌克兰</option>
<option value="UY">乌拉圭</option>
<option value="UZ">乌兹别克斯坦</option>
<option value="ES">西班牙</option>
<option value="GR">希腊</option>
<option value="HK">香港</option>
<option value="SG">新加坡</option>
<option value="NZ">新西兰</option>
<option value="HU">匈牙利</option>
<option value="JM">牙买加</option>
<option value="AM">亚美尼亚</option>
<option value="YE">也门</option>
<option value="IL">以色列</option>
<option value="IT">意大利</option>
<option value="IN">印度</option>
<option value="ID">印度尼西亚</option>
<option value="GB">英国</option>
<option value="JO">约旦</option>
<option value="VN">越南</option>
<option value="TD">乍得</option>
<option value="CL">智利</option>
<option selected="selected" value="CN">中国</option>
<option value="MO">中国澳门特别行政区</option>
<option value="FJ">斐济</option></select>
                            </div>
                            <input type="submit" value="更改" name="2.0.1.1.3.0.7.11.3.1.0.5.11.1.0.5.7" />
                            <p class="disclaimer">注意：更改国家或地区后，页面上的语言可能也会改变。</p>


                        </div>
                    </div>
                </div>



                <div class="billing-information">
                    <div class="account-fields">





                        <div class="payment-info cup">
                            <div class="formset credit-card logos-v2 cup">
                                <label>付款方式</label>
                                <input type="hidden" name="credit-card-type" value="" />
                                <input name="sp" type="hidden" value='' />
                                <input name="res" type="hidden" value='' />
                                <ul class="credit-card-picker" role="radiogroup">

                                    <li role="presentation"><input role="radio" aria-checked="false" selection="None" class="upcc" id="cc_upcc" style="background-image: url(https://s.mzstatic.com/images/creditcards/-dsi-/cc_upcc.png);" title="UnionPay" type="submit" value="UnionPay" name="UPCC" /></li>

                                    <li role="presentation"><input role="radio" aria-checked="false" selection="None" class="payease" id="cc_payease" style="background-image: url(https://s.mzstatic.com/images/creditcards/-dsi-/cc_peas.png);" type="submit" value="" name="PayEase" /></li>

                                    <li role="presentation"><input role="radio" aria-checked="false" selection="None" class="visa" id="cc_visa" style="background-image: url(https://s.mzstatic.com/images/creditcards/-dsi-/cc_visa.png);" title="Visa" type="submit" value="Visa" name="Visa" /></li>

                                    <li role="presentation"><input role="radio" aria-checked="false" selection="None" class="mastercard" id="cc_mastercard" style="background-image: url(https://s.mzstatic.com/images/creditcards/-dsi-/cc_mast.png);" title="MasterCard" type="submit" value="MasterCard" name="MasterCard" /></li>

                                    <li role="presentation"><input role="radio" aria-checked="false" selection="None" class="amex" id="cc_amex" style="background-image: url(https://s.mzstatic.com/images/creditcards/-dsi-/cc_amex.png);" title="Amex" type="submit" value="Amex" name="Amex" /></li>

                                    <li role="presentation"><input role="radio" aria-checked="true" selection="None" class="none selected" id="cc_none" title="已选择 无" type="submit" value="无" name="None" /></li>

                                </ul>



                            </div>





                            <div class="cup-msg-block">
                                <p class="instruction">你的银行必须开通你的卡才能网上购物。</p>
                            </div>



                            <div class="formset optional-section hidden" id="cup_fields">



                                <label>移动电话</label>
                                <div>
                                    <div class="formset text mobile-phone stacked"><input viewName="mobilePhone" placeholder="12345678901" class="text mobile-phone optional" id="mobile-phone" title="手机
" type="text" name="2.0.1.1.3.0.7.11.3.1.0.5.17.0.0.0.17.5" /></div>
                                </div>

                                <div class="cup-mobile-msg-block">
                                    <p class="instruction">输入银行预留手机号。</p>
                                </div>


                            </div>
                            <input class="text text optional" id="card_type_id" type="hidden" value="0" name="card-type-id" />


                            <div style="display:none" id="getCardTypeUrl">https://play.itunes.apple.com/WebObjects/MZPlay.woa/wa/getCardTypeSrv</div>


                            <div class="paypal-info">
                                要确认您的 PayPal 帐户，并确定您愿意将其用于在 iTunes Store 购物，请点击“继续”。
                                <a id="paypallink" href="https://www.paypal.com/cgi-bin/webscr?business=paypalgroup%40group.apple.com&mp_max_edit=0&mp_desc=iTunes+Store+%E8%B4%AD%E4%B9%B0%E9%A1%B9%E7%9B%AE%E3%80%82&workflowID=16517&mp_max=500&mp_test_amount=20.00&mp_custom=sessionID%3DXXXX%26workflowID%3D16517&return=https%3A%2F%2Fbuy.itunes.apple.com%2FWebObjects%2FMZFinance.woa%2F-1%2Fwa%2Fcom.apple.jingle.app.finance.DirectAction%2FpayPalMIPReturn&mp_pay_type=i&country=CN&lc=zh&cmd=_xclick-merchant&mp_max_min=700&mp_fs_country=US"></a>
                            </div>
                        </div>





                        <div class="code-redemption">
                            <hr/>



                            <p class="instruction">要兑现代码或礼品券，请在此处输入。</p>



                            <div class="boxed inset">
                                <div class="formset text before field"><label class="text before validated optional" for="codeRedemptionField"><span>代码</span></label><input autocapitalize="off" autocomplete="off" autocorrect="off" placeholder="输入代码" class="text before validated optional" id="codeRedemptionField" title="输入代码" type="text" name="2.0.1.1.3.0.7.11.3.1.0.5.21.1.3.1.3.2.3" /></div>
                            </div>





                        </div>




                        <hr/>
                        <div class="address-information">




                            <p class="instruction">帐单寄送地址</p>





                            <div class="address-block">


                                <div class="name-block">

                                    <div class="formset text"><input viewName="lastName" placeholder="姓" class="text" id="lastFirstName" title="姓" type="text" name="2.0.1.1.3.0.7.11.3.1.0.5.23.5.1.3.1" /></div>
                                    <div class="formset text"><input viewName="firstName" placeholder="名" class="text validated required" id="firstName" title="名" type="text" name="2.0.1.1.3.0.7.11.3.1.0.5.23.5.1.3.3" /></div>





                                </div>





                                <div class="street-block">
                                    <div class="formset text"><input viewName="street1" placeholder="街名和门牌号" class="text" id="street1" title="街名和门牌号" type="text" name="2.0.1.1.3.0.7.11.3.1.0.5.23.5.7.13.1" /></div>
                                    <div class="formset text"><input viewName="street2" placeholder="楼号、单元号、房间号" class="text optional" id="street2" title="楼号、单元号、房间号" type="text" name="2.0.1.1.3.0.7.11.3.1.0.5.23.5.7.13.3" /></div>
                                    <div class="formset text"><input viewName="street3" placeholder="街" class="text optional" id="street3" title="街" type="text" name="2.0.1.1.3.0.7.11.3.1.0.5.23.5.7.13.5" /></div>
                                </div>

                                <div class="address-group">
                                    <div class="formset text"><input viewName="city" placeholder="市级行政区" class="text optional" id="city" title="市级行政区" type="text" name="2.0.1.1.3.0.7.11.3.1.0.5.23.5.7.13.7" /></div>
                                </div>
                                <div class="address-group">
                                    <div class="formset text postalCode"><input viewName="postalCode" placeholder="邮编" class="text optional" id="postalcode" title="邮编" type="text" name="2.0.1.1.3.0.7.11.3.1.0.5.23.5.7.13.9" /></div>
                                    <div class="formset select"><select viewName="state" placeholder="选择省份" class="select optional" id="state" title="选择省份" name="state"><option value="WONoSelectionString">省份</option>
<option value="0">上海</option>
<option value="1">云南</option>
<option value="2">内蒙古</option>
<option value="3">北京</option>
<option value="4">吉林</option>
<option value="5">四川</option>
<option value="6">天津</option>
<option value="7">宁夏</option>
<option value="8">安徽</option>
<option value="9">山东</option>
<option value="10">山西</option>
<option value="11">广东</option>
<option value="12">广西</option>
<option value="13">新疆</option>
<option value="14">江苏</option>
<option value="15">江西</option>
<option value="16">河北</option>
<option value="17">河南</option>
<option value="18">浙江</option>
<option value="19">海南</option>
<option value="20">湖北</option>
<option value="21">湖南</option>
<option value="22">甘肃</option>
<option value="23">福建</option>
<option value="24">西藏</option>
<option value="25">贵州</option>
<option value="26">辽宁</option>
<option value="27">重庆</option>
<option value="28">陕西</option>
<option value="29">青海</option>
<option value="30">黑龙江</option></select></div>
                                    <span class="country-name">中国</span>
                                </div>
                                <div class="phone-block">

                                    <div class="formset text phone"><input viewName="phone1Number" placeholder="电话" class="text phone" id="phone1Number" title="电话" type="tel" name="2.0.1.1.3.0.7.11.3.1.0.5.23.5.7.13.17" /></div>

                                </div>





                            </div>





                        </div>





                        <div class="captcha-container">
                            <script id="nucaptcha-template" type="text/x-nucaptcha-template">
                                ${PlayerStart}
                                <p class="instruction captcha-instruction">输入动态字符：</p>
                                <div class="captcha-controls">
                                    ${InputAnswer}
                                    <div class="pill" role="radio-group">
                                        <button id="player-mode-video" class="cmd toggle active text-captcha" role="radio" aria-label="取得文本验证码" aria-checked="true" type="button"><span>取得文本验证码</span></button>
                                        <button id="player-mode-audio" class="cmd toggle audio-captcha" role="radio" aria-label="取得音频验证码" aria-checked="false" type="button"><span>取得音频验证码</span></button>
                                    </div>
                                    <button id="new-challenge" class="cmd new-captcha" type="button"><span>取得新验证码</span></button>
                                </div>
                                <div class="media-wrapper">
                                    ${Media}
                                    <button class="replay-button" aria-label="重播验证码" type="button"><span>重播验证码</span></button>
                                </div>
                                ${PlayerEnd}
                            </script>
                            <div id="ndwc"><noscript><input type="hidden" id="ndpd-s-ns" name="ndpd-s-ns" value="1.w-855182.1.2.7GT6A1_6NpLUtdhFQF8KEw,,.Ons7F8ExxTeO2PYwGY90Qf086f9sKpYb8QleTRQ87r-edMFZqULXM6sh81zTX8356xneENnFZ8Hwzndqn_XtNCO_v-JtnCDPR-6Z1kWi1Clruod6pf8TkiO1GYVrNm-H08UPPdv5Zbc7ZAzGhb3QQbfv0mgZ5iEI3W4ixGQQvpyhZz0j0FRg8SYI9OTdTEvFQ5WynthC0veYspBrwdhqDmIOUJ_ccFs6rDnDaWtsjKycT3lEa8TVRmOhvg2OJOaeobnKdBEsJjKzmS_1lQkl80RKaOW6jYFTn3CZgMLtP53qg3ei1zimtU6WT2yR-wH6kVGC9YC9qSCDCnWYbSboew,,"><input type="hidden" id="ndpd-f-ns" name="ndpd-f-ns" value="1.w-855182.1.2.Cuvu2SaSndTN3B3Cvyr79Q,,.fmPv0yeYTcoN6QQFs9FL_QXigGBrQIVMjQz3xhQjb3yNIu7UTu9WAxLAuNLlZg-uhVFIIXmjrj3j3DZfuDYEVecs2fpQL4VLVViJbAe-TAIBa5Y0DcjZt9i2ZfJmqTw0oA1qgUwkkPOUfzwshz1TEiG728bgojIdFMmyrJ_g2RFr0eerMEcabcdl7Vj3sRfypHOJaM5btWQ5JDDfA-Rol21TUgn8FWtRIXK0GXP_Ykh6VVHNg-ad4-x4EsgqWQ98KA-Q-Xv6eW1-1G-6R-MB8WGUReMtgV38BCTmAoz1OJ9myTA6ILXmqTKLIZEs73lQXDPZ9mOPnpdXWqbSgyjb6g,,"><input type="hidden" id="ndpd-fm-ns" name="ndpd-fm-ns" value="1.w-855182.1.2.V02C5dfySZ9YRj_KK6JeXQ,,.j_0v9peSPIHF-W1qGp5ESgac5hpkuSS7YTBAQ2Xd_mXgmub1WLhodleCuY1ZSVvPhoK6P53xfZoi0LJiADNEpNvNw6iaypVKS40Ze34VO1UEOYvOj1haQcgRsy8wmnRD-jlLAiqrvXQd-AqG1Eqst5JgxP5EjoXHyn732I5jpB6_s_8f3JIz3IImd78u1SlJPYOj2FLYPXCesi-4p1lwhJrDTJEcstC6H-J5NCsvw2T1INPwE8J0gbzsb4W_jKqW2Dm8cOfB6EM-t6DR9lcVcpQKk1k3KyRNtgFKggEdgFpjV6Bk5_HI-akw2k0TxOirgwI8zucUFESz3DFeGhL0iA,,"><input type="hidden" id="ndpd-w-ns" name="ndpd-w-ns" value="1.w-855182.1.2.N4D4th2mzU7DZFe4GOCXSQ,,.DR2hV8-NJlKDoVrGV0VB_KksOn06w7kADkjIfzO-lIC8zHUoS4h-_dn1G-5bu44gF9NjitGEOqRMyc-CKm8YKo18JJV7xV33v8d11v1o3926J-j_tlPGyYv4AcI4UfaJXq1YZ3hzmhmjUKB-GB9oFXjcX7kTUsZQoUUKMH_D6wwvhM_JqKdNWqO3f4tPZWAwEHfi8ZHE5d_3Sf6EohA6L2ZpExB1gFR93MYy8naPudp-_3L9Sz2Y92Guw0CdfCufUheEpWhRDhmi87YMvZmnFu7O2tGR3eI8WuuKtRne-dI-BHgqWBVle0ZomRKDMw8bvO_XPi8MbI913g4E3Ycwx4PIuTXqBWDdjSdb-JWDYI4,"></noscript><input type="hidden" id="ndpd-s" name="ndpd-s" value="1.w-855182.1.2.7GT6A1_6NpLUtdhFQF8KEw,,.Ons7F8ExxTeO2PYwGY90Qf086f9sKpYb8QleTRQ87r-edMFZqULXM6sh81zTX8356xneENnFZ8Hwzndqn_XtNCO_v-JtnCDPR-6Z1kWi1Clruod6pf8TkiO1GYVrNm-H08UPPdv5Zbc7ZAzGhb3QQbfv0mgZ5iEI3W4ixGQQvpyhZz0j0FRg8SYI9OTdTEvFQ5WynthC0veYspBrwdhqDmIOUJ_ccFs6rDnDaWtsjKycT3lEa8TVRmOhvg2OJOaeobnKdBEsJjKzmS_1lQkl80RKaOW6jYFTn3CZgMLtP53qg3ei1zimtU6WT2yR-wH6kVGC9YC9qSCDCnWYbSboew,,"><input type="hidden" id="ndpd-f" name="ndpd-f" value="1.w-855182.1.2.Cuvu2SaSndTN3B3Cvyr79Q,,.fmPv0yeYTcoN6QQFs9FL_QXigGBrQIVMjQz3xhQjb3yNIu7UTu9WAxLAuNLlZg-uhVFIIXmjrj3j3DZfuDYEVecs2fpQL4VLVViJbAe-TAIBa5Y0DcjZt9i2ZfJmqTw0oA1qgUwkkPOUfzwshz1TEiG728bgojIdFMmyrJ_g2RFr0eerMEcabcdl7Vj3sRfypHOJaM5btWQ5JDDfA-Rol21TUgn8FWtRIXK0GXP_Ykh6VVHNg-ad4-x4EsgqWQ98KA-Q-Xv6eW1-1G-6R-MB8WGUReMtgV38BCTmAoz1OJ9myTA6ILXmqTKLIZEs73lQXDPZ9mOPnpdXWqbSgyjb6g,,"><input type="hidden" id="ndpd-fm" name="ndpd-fm" value="1.w-855182.1.2.V02C5dfySZ9YRj_KK6JeXQ,,.j_0v9peSPIHF-W1qGp5ESgac5hpkuSS7YTBAQ2Xd_mXgmub1WLhodleCuY1ZSVvPhoK6P53xfZoi0LJiADNEpNvNw6iaypVKS40Ze34VO1UEOYvOj1haQcgRsy8wmnRD-jlLAiqrvXQd-AqG1Eqst5JgxP5EjoXHyn732I5jpB6_s_8f3JIz3IImd78u1SlJPYOj2FLYPXCesi-4p1lwhJrDTJEcstC6H-J5NCsvw2T1INPwE8J0gbzsb4W_jKqW2Dm8cOfB6EM-t6DR9lcVcpQKk1k3KyRNtgFKggEdgFpjV6Bk5_HI-akw2k0TxOirgwI8zucUFESz3DFeGhL0iA,,"><input type="hidden" id="ndpd-w" name="ndpd-w" value="1.w-855182.1.2.N4D4th2mzU7DZFe4GOCXSQ,,.DR2hV8-NJlKDoVrGV0VB_KksOn06w7kADkjIfzO-lIC8zHUoS4h-_dn1G-5bu44gF9NjitGEOqRMyc-CKm8YKo18JJV7xV33v8d11v1o3926J-j_tlPGyYv4AcI4UfaJXq1YZ3hzmhmjUKB-GB9oFXjcX7kTUsZQoUUKMH_D6wwvhM_JqKdNWqO3f4tPZWAwEHfi8ZHE5d_3Sf6EohA6L2ZpExB1gFR93MYy8naPudp-_3L9Sz2Y92Guw0CdfCufUheEpWhRDhmi87YMvZmnFu7O2tGR3eI8WuuKtRne-dI-BHgqWBVle0ZomRKDMw8bvO_XPi8MbI913g4E3Ycwx4PIuTXqBWDdjSdb-JWDYI4,"><input type="hidden" id="ndpd-ipr" name="ndpd-ipr" value=""><input type="hidden" id="ndpd-di" name="ndpd-di" value="p"><input type="hidden" id="ndpd-bi" name="ndpd-bi" value="p"><input type="hidden" id="ndpd-wk" name="ndpd-wk" value="p"><input type="hidden" id="ndpd-vk" name="ndpd-vk" value="6622"></div>
                            <script type="text/javascript">
                                var nsqpd, nsqpdp, nspdbbpddp, nsdqq = {},
                                    nsdqqbdqqd = {},
                                    nsdqbp = -1,
                                    nsqpbpd = -1,
                                    nsdqqb = [],
                                    nsqpbpdqqd = [],
                                    nsdbpdbqd = "fspm",
                                    nspdbbp = null;

                                function ndwti(a) {
                                    nsdqbp = nsppbdqq();
                                    typeof a === "string" && (a = nspdqpppqp(a));
                                    nsqpd = a.did;
                                    nsqpdp = a.fff;
                                    nspdbbpddp = a.ffft;
                                    nsdqq = a.wmd;
                                    nsdqqbdqqd = a.fd;
                                    ndovFormMode = a.ffmm;
                                    ndovSingleName = a.fsss;
                                    ndovSingleDKV = a.ddkv;
                                    ndovSingleDF = a.dddf;
                                    ndovIPRDelimMode = a.ppns;
                                    ndovIPRDelimModeALT = a.ppmm;
                                    ndovIPRDelimALT = a.ppdd;
                                    ndovIPRDelimStand = a.ppds;
                                    ndovWidgetKeyOn = a.wwwe;
                                    if (typeof a.mp !== "undefined") nspdbbp = a.mp;
                                    nsdbpdbqdp();
                                    for (var c = 0; c < nsqpbpdqqd.length; c++)(0, nsqpbpdqqd[c][1])(a, nsdqq[nsqpbpdqqd[c][0]]);
                                    nsdbpdbqdp();
                                    nsqpbpd = nsppbdqq();
                                    nsdbpdbqdp()
                                }

                                function ndwtr() {
                                    for (var a = 0; a < nsqpbpdqqd.length; a++)
                                        if (nsqpbpdqqd[a].length >= 3 && typeof nsqpbpdqqd[a][2] !== "undefined")(0, nsqpbpdqqd[a][2])();
                                    nsdbpdbqdp()
                                }

                                function nsdqbpbd(a) {
                                    return nsqpdp.replace(nspdbbpddp, a)
                                }

                                function nspdb(a) {
                                    nsdqqb.push(a);
                                    nspqdqqp(nsdqbpbd("jse"), nsdqbpbdb.stringify(nsdqqb))
                                }

                                function nspdbbpd(a, c) {
                                    nsdqqbdqqd[a] = c
                                }

                                function nsdbpdbqdp() {
                                    field = "";
                                    for (var a in nsdqqbdqqd) ndovFormMode === nsdbpdbqd ? field += nsdqbpbd(a) + ndovSingleDKV + nsdqqbdqqd[a] + ndovSingleDF : nspqdqqp(nsdqbpbd(a), nsdqqbdqqd[a]);
                                    ndovFormMode === nsdbpdbqd && (field.substring(field.length - ndovSingleDF.length, field.length) === ndovSingleDF && (field = field.substring(0, field.length - ndovSingleDF.length)), nspqdqqp(ndovSingleName, field))
                                }

                                function nspdppddd(a, c, b) {
                                    nsqpbpdqqd.push([a, c, b])
                                }

                                function nsqddqb(a) {
                                    return a
                                }

                                function nspdqpppqp(a) {
                                    return a
                                }

                                function nspqdqqp(a, c) {
                                    var b = [""];
                                    nspdbbp !== null && (b = nspdbbp);
                                    for (var e = 0; e < b.length; e++) {
                                        var d = b[e];
                                        d != "" && (d = "-" + d);
                                        var f = nsbbbd(a + d);
                                        if (f !== null) f.value = c;
                                        else {
                                            var h = nsbbbd(nsqpd + d),
                                                f = document.createElement("input");
                                            f.id = a + d;
                                            f.name = a;
                                            f.value = c;
                                            f.type = "hidden";
                                            h.appendChild(f)
                                        }
                                    }
                                }

                                function nsppbdqq() {
                                    return parseInt((new Date).getTime() / 1E3)
                                }

                                function nsbpdqb() {
                                    return parseInt((new Date).getTime())
                                }

                                function nsqpbqdqq(a) {
                                    var c = document.createElement("script");
                                    c.setAttribute("type", "text/JavaScript");
                                    c.setAttribute("src", a);
                                    document.getElementsByTagName("head")[0].appendChild(c)
                                }

                                function nsbbbd(a) {
                                    var c = null;
                                    document.getElementById ? c = document.getElementById(a) : document.all && (c = document.all[a]);
                                    return c
                                }

                                function nsbbpddbp(a, c) {
                                    arReturn = [];
                                    typeof a.getElementsByTagName !== "undefined" && (arReturn = a.getElementsByTagName(c));
                                    return arReturn
                                }
                                var nsdqbpbdb;
                                nsdqbpbdb || (nsdqbpbdb = {});
                                (function() {
                                    function a(a) {
                                        return a < 10 ? "0" + a : a
                                    }

                                    function c(a) {
                                        d.lastIndex = 0;
                                        return d.test(a) ? '"' + a.replace(d, function(a) {
                                            var b = n[a];
                                            return typeof b === "string" ? b : "\\u" + ("0000" + a.charCodeAt(0).toString(16)).slice(-4)
                                        }) + '"' : '"' + a + '"'
                                    }

                                    function b(a, d) {
                                        var e, m, k, o, p = f,
                                            j, g = d[a];
                                        g && typeof g === "object" && typeof g.toNDJSON === "function" && (g = g.toNDJSON(a));
                                        typeof i === "function" && (g = i.call(d, a, g));
                                        switch (typeof g) {
                                            case "string":
                                                return c(g);
                                            case "number":
                                                return isFinite(g) ? String(g) : "null";
                                            case "boolean":
                                            case "null":
                                                return String(g);
                                            case "object":
                                                if (!g) return "null";
                                                f += h;
                                                j = [];
                                                if (Object.prototype.toString.apply(g) === "[object Array]") {
                                                    o = g.length;
                                                    for (e = 0; e < o; e += 1) j[e] = b(e, g) || "null";
                                                    k = j.length === 0 ? "[]" : f ? "[\n" + f + j.join(",\n" + f) + "\n" + p + "]" : "[" + j.join(",") + "]";
                                                    f = p;
                                                    return k
                                                }
                                                if (i && typeof i === "object") {
                                                    o = i.length;
                                                    for (e = 0; e < o; e += 1) typeof i[e] === "string" && (m = i[e], (k = b(m, g)) && j.push(c(m) + (f ? ": " : ":") + k))
                                                } else
                                                    for (m in g) Object.prototype.hasOwnProperty.call(g, m) && (k = b(m, g)) && j.push(c(m) + (f ? ": " : ":") + k);
                                                k = j.length === 0 ? "{}" : f ? "{\n" + f + j.join(",\n" +
                                                    f) + "\n" + p + "}" : "{" + j.join(",") + "}";
                                                f = p;
                                                return k
                                        }
                                    }
                                    if (typeof Date.prototype.toNDJSON !== "function") Date.prototype.toNDJSON = function() {
                                        return isFinite(this.valueOf()) ? this.getUTCFullYear() + "-" + a(this.getUTCMonth() + 1) + "-" + a(this.getUTCDate()) + "T" + a(this.getUTCHours()) + ":" + a(this.getUTCMinutes()) + ":" + a(this.getUTCSeconds()) + "Z" : null
                                    }, String.prototype.toNDJSON = Number.prototype.toNDJSON = Boolean.prototype.toNDJSON = function() {
                                        return this.valueOf()
                                    };
                                    var e = /[\u0000\u00ad\u0600-\u0604\u070f\u17b4\u17b5\u200c-\u200f\u2028-\u202f\u2060-\u206f\ufeff\ufff0-\uffff]/g,
                                        d = /[\\\"\x00-\x1f\x7f-\x9f\u00ad\u0600-\u0604\u070f\u17b4\u17b5\u200c-\u200f\u2028-\u202f\u2060-\u206f\ufeff\ufff0-\uffff]/g,
                                        f, h, n = {
                                            "\u0008": "\\b",
                                            "\t": "\\t",
                                            "\n": "\\n",
                                            "\u000c": "\\f",
                                            "\r": "\\r",
                                            '"': '\\"',
                                            "\\": "\\\\"
                                        },
                                        i;
                                    if (typeof nsdqbpbdb.stringify !== "function") nsdqbpbdb.stringify = function(a, c, d) {
                                        var e;
                                        h = f = "";
                                        if (typeof d === "number")
                                            for (e = 0; e < d; e += 1) h += " ";
                                        else typeof d === "string" && (h = d);
                                        if ((i = c) && typeof c !== "function" && (typeof c !== "object" || typeof c.length !== "number")) throw Error("nsdqbpbdb.stringify");
                                        return b("", {
                                            "": a
                                        })
                                    };
                                    if (typeof nsdqbpbdb.parse !== "function") nsdqbpbdb.parse = function(a, b) {
                                        function c(a, d) {
                                            var e, i, f = a[d];
                                            if (f && typeof f === "object")
                                                for (e in f) Object.prototype.hasOwnProperty.call(f, e) && (i = c(f, e), i !== void 0 ? f[e] = i : delete f[e]);
                                            return b.call(a, d, f)
                                        }
                                        var d, a = String(a);
                                        e.lastIndex = 0;
                                        e.test(a) && (a = a.replace(e, function(a) {
                                            return "\\u" + ("0000" + a.charCodeAt(0).toString(16)).slice(-4)
                                        }));
                                        if (/^[\],:{}\s]*$/.test(a.replace(/\\(?:["\\\/bfnrt]|u[0-9a-fA-F]{4})/g, "@").replace(/"[^"\\\n\r]*"|true|false|null|-?\d+(?:\.\d*)?(?:[eE][+\-]?\d+)?/g,
                                                "]").replace(/(?:^|:|,)(?:\s*\[)+/g, ""))) return d = eval("(" + a + ")"), typeof b === "function" ? c({
                                            "": d
                                        }, "") : d;
                                        throw new SyntaxError("nsdqbpbdb.parse");
                                    }
                                })();

                                function nsqddqbdb(a, c) {
                                    if (ndovWidgetKeyOn) {
                                        var b = Math.floor(Math.random() * 1E6) + 1E3;
                                        nspdbbpd("wkr", b);
                                        nsqpbqdqq(c.r + "?r=" + b + "&wt=" + c.w);
                                        nspdbbpd("wk", "p")
                                    }
                                }
                                nspdppddd("wk", nsqddqbdb);

                                function ndwtw(a) {
                                    ndovWidgetKeyOn && (typeof a === "string" && (a = nspdqpppqp(a)), nspdbbpd("wk", a.wk), nsdbpdbqdp())
                                }
                                var nspdppdd = 0,
                                    nspqqqbd = null,
                                    nsqpbp = false,
                                    nspdppdddp = null,
                                    nspqqqb = /^(text|password|email|url|search|tel)$/i;

                                function nsbbpdd(a, c) {
                                    nspqqqbd = c.fm;
                                    nsqpbp = "lm" in c && c.lm;
                                    nsbpd(nspqqqbd, nsqpbp);
                                    nspdbbpd("ipr", "");
                                    nsbbb()
                                }

                                function nspdqp() {}
                                nspdppddd("ipr", nsbbpdd, nspdqp);

                                function nspdqpppq(a, c, b) {
                                    window.addEventListener ? a.addEventListener(c, b, false) : a.attachEvent && a.attachEvent("on" + c, b)
                                }

                                function nsbbb() {
                                    var a = nsppbdqqpb();
                                    nsqpbqd(nsdqbpbdbq, [a])
                                }

                                function nsppbdqqpb() {
                                    for (var a = [], c = 0; c < nspdppdddp.length; c++) {
                                        var b = nspdppdddp[c];
                                        b.type && b.type.match(nspqqqb) && (a.push(b.id), a.push(b.value.length))
                                    }
                                    return a.join(",")
                                }
                                var nsdbpd = "a",
                                    nspdp = "b",
                                    nsqpbpdqq = "c",
                                    nsqpdpqqbb = "d",
                                    nsdqb = "e",
                                    nspqq = "f";

                                function nsbbbdbpqp(a) {
                                    if (!a) a = window.event;
                                    var c = null;
                                    if (a.target) c = a.target;
                                    else if (a.srcElement) c = a.srcElement;
                                    if (c.nodeType == 3) c = c.parentNode;
                                    var b;
                                    if (a.keyCode) b = a.keyCode;
                                    else if (a.which) b = a.which;
                                    var e = false;
                                    a.which ? e = a.which == 3 : a.button && (e = a.button == 2);
                                    var d = 0,
                                        f = 0;
                                    if (a.pageX || a.pageY) d = a.pageX, f = a.pageY;
                                    else if (a.clientX || a.clientY) d = a.clientX + document.body.scrollLeft + document.documentElement.scrollLeft, f = a.clientY + document.body.scrollTop + document.documentElement.scrollTop;
                                    var h = {};
                                    h[nsdbpd] =
                                        a;
                                    h[nspdp] = c;
                                    h[nsqpbpdqq] = e;
                                    h[nsqpdpqqbb] = b;
                                    h[nsdqb] = d;
                                    h[nspqq] = f;
                                    return h
                                }

                                function nsppbdq(a) {
                                    return !nspdp in a ? null : typeof a[nspdp].id === "string" && a[nspdp].id !== "" ? a[nspdp].id : a[nspdp].name
                                }

                                function nsbpd(a) {
                                    for (var c = null, b = arguments.length > 1 ? arguments[1] : FALSE, e = [], d = document.documentElement; d.childNodes.length && d.lastChild.nodeType === 1;) d = d.lastChild;
                                    for (; typeof d.parentNode !== "undefined" && null === d.tagName.match(/(form|html)/i);) d = d.parentNode;
                                    d.tagName.match(/(form|html)/i) && (c = d, e = nsbbpddbp(c, "input"));
                                    if (b && c !== null) {
                                        a = [];
                                        for (b = 0; b < e.length; b++) {
                                            var f = e[b];
                                            f.type && f.type.match(nspqqqb) && a.push(f)
                                        }
                                    } else {
                                        for (var h = [], b = 0; b < a.length; b++) {
                                            d = nsbbbd(a[b]);
                                            if (null === d && null !== c)
                                                for (var n =
                                                        0; n < e.length; n++) f = e[n], f.type && f.type.match(nspqqqb) && f.name && f.name === a[b] && (d = f);
                                            null !== d && h.push(d)
                                        }
                                        a = h
                                    }
                                    nspdppdddp = a;
                                    for (b = 0; b < a.length; b++)
                                        if (d = a[b], null !== d) {
                                            if (c === null) {
                                                e = d.parentNode;
                                                for (f = 0; f < 8; f++) {
                                                    if (e === null || e.nodeName.match(/form/i)) {
                                                        c = e;
                                                        break
                                                    }
                                                    e = e.parentNode
                                                }
                                            }
                                            d.nodeName.match(/input/i) && (nspdqpppq(d, "keydown", function(a) {
                                                nsbbbdbpqp(a);
                                                nsqpbqd(nspdbb, [])
                                            }), nspdqpppq(d, "focus", function(a) {
                                                a = nsbbbdbpqp(a);
                                                nsqpbqd(nsdqqbdq, [nspdp in a && typeof a[nspdp].value !== "undefined" ? a[nspdp].value.length :
                                                    null, nsppbdq(a)
                                                ]);
                                                nsqpbqd(nspdbbpdd, [nsppbdq(a)])
                                            }), nspdqpppq(d, "blur", function(a) {
                                                a = nsbbbdbpqp(a);
                                                nsqpbqd(nsdbp, [nsppbdq(a)])
                                            }))
                                        }
                                    nspdqpppq(document, "click", function(a) {
                                        a = nsbbbdbpqp(a);
                                        nsqpbqd(nsqpdpq, [a[nsdqb], a[nspqq], nsppbdq(a)])
                                    });
                                    nspdqpppq(document, "touchstart", function(a) {
                                        a = nsbbbdbpqp(a);
                                        a.event && a.event.touches && a.event.touches[0] && typeof a.event.touches[0].pageX !== "undefined" ? nsqpbqd(nspdpp, [a.event.touches[0].pageX, a.event.touches[0].pageY, nsppbdq(a)]) : nsqpbqd(nspdpp, [-1, -1, nsppbdq(a)])
                                    });
                                    nspdqpppq(document, "mousemove", function(a) {
                                        if (nsppbdqq() < nspdppdd) return false;
                                        nspdppdd = nsppbdqq() + 5;
                                        a = nsbbbdbpqp(a);
                                        nsqpbqd(nspdppd, [a[nsdqb], a[nspqq], nsppbdq(a)])
                                    });
                                    null !== c && nspdqpppq(c, "submit", function(a) {
                                        a = nsbbbdbpqp(a);
                                        nsqpbqd(nsqpdpqq, [a[nsdqb], a[nspqq], c.id])
                                    })
                                }
                                var nspdbbpdd = "ff",
                                    nsdbp = "fb",
                                    nspdbb = "kd",
                                    nsdbpdbq = "ku",
                                    nspdppd = "mm",
                                    nsqpdpq = "mc",
                                    nsqpdpqqb = "ac",
                                    nspdpp = "te",
                                    nsqpdpqq = "fs",
                                    nspqqq = "sp",
                                    nsdqqbdq = "kk",
                                    nsdqbpbdbq = "st",
                                    nsdqqbdqq = 1,
                                    nsdqqbd = 1,
                                    nsqpb = null,
                                    nsdqbpb = null,
                                    nspqqqbdqb = null,
                                    nsdbpdb = null,
                                    nspqqqbdq = "",
                                    nsqpbpdq = "";

                                function nsqpbqd(a, c) {
                                    var b = nsbpdqb();
                                    nsdqbpb == null && (nsdbpdb = nsdqbpb = nsqpb = nsbpdqb(), nsqpbqdq("ncip", b, [nsppbdqq(), nsdqqbdqq, nsdqqbd]));
                                    nsqpbqdq(a, b, c);
                                    b - nsdbpdb >= 15E3 && (nsqpbqdq("ts", b, [b - nsqpb]), nsdbpdb = b);
                                    switch (a) {
                                        case nsdbp:
                                        case nsdbp:
                                        case nsqpdpq:
                                        case nsqpdpqq:
                                            nspqdqqpbd(b);
                                            break;
                                        default:
                                            b - nspqqqbdqb > 2E3 && nspqdqqpbd(b)
                                    }
                                }

                                function nspqdqqpbd(a) {
                                    nspqqqbdqb = a;
                                    a = "";
                                    nspqqqbdq !== "" && (nsqpbpdq += nspqqqbdq, nspqqqbdq = "", a = nsqpbpdq, ndovIPRDelimMode === ndovIPRDelimModeALT && (a = a.replace(RegExp(ndovIPRDelimStand, "g"), ndovIPRDelimALT)), nspdbbpd("ipr", a), nsdbpdbqdp())
                                }

                                function nsqpbqdq(a, c, b) {
                                    var e = c - nsdqbpb;
                                    nsdqqbd > 1 && (e = Math.round(e / nsdqqbd));
                                    a = a + "," + e.toString(16);
                                    if (b != null && b.length > 0) {
                                        for (var e = [], d = 0; d < b.length; d++) typeof b[d] === "number" ? e.push(Math.round(b[d]).toString(16)) : null != b[d] && e.push(b[d].toString());
                                        a = a + "," + e.join(",")
                                    }
                                    nspqqqbdq = nspqqqbdq + a + ";";
                                    nsdqbpb = c
                                }

                                function nsbbpddbpd() {
                                    nspdbbpd("di", nsppbdqqp());
                                    nspdbbpd("bi", nsbbpddb())
                                }
                                nspdppddd("di", nsbbpddbpd);

                                function nspdqpp() {
                                    var a = [];
                                    a.push(nspqdqqpb());
                                    a.push(nsbpdqbbdd());
                                    a.push(nsqddqbd());
                                    a.push(nsbbpd());
                                    for (var c = "DI", b = 0; b < a.length; b++) c += "." + a[b];
                                    return c
                                }

                                function nsbbpddb() {
                                    var a = [];
                                    a.push(nsbpdqbbdd());
                                    a.push(nsqddqbd());
                                    a.push(nsppbd());
                                    for (var c = "b1", b = 0; b < a.length; b++) c += "." + a[b];
                                    return c
                                }

                                function nsppbdqqp() {
                                    return "d1-" + nspqdqq(nspdqpp())
                                }

                                function nspqdqqpb() {
                                    var a = "NotAvail";
                                    if (typeof navigator !== "undefined" && typeof navigator.userAgent !== "undefined") a = navigator.userAgent, a = a.replace(/([0-9]+\.[0-9]+)\.[0-9]+\.[0-9]+/g, "$1").replace(/([0-9]+\.[0-9]+)\.[0-9]+/g, "$1"), a = a.replace(/([0-9]+_[0-9]+)_[0-9]+_[0-9]+/g, "$1").replace(/([0-9]+_[0-9]+)_[0-9]+/g, "$1");
                                    return a
                                }

                                function nsppbd() {
                                    return window.navigator.userLanguage || window.navigator.language || window.navigator.browserLanguage
                                }

                                function nsbpdqbbdd() {
                                    var a = "";
                                    typeof window.screen !== "undefined" && (typeof(window.screen.width !== "undefined") && typeof(window.screen.height !== "undefined") && (a += window.screen.width + "x" + window.screen.height), typeof(window.screen.availWidth !== "undefined") && typeof(window.screen.availHeight !== "undefined") && (a += " " + window.screen.availWidth + "x" + window.screen.availHeight), typeof(window.screen.colorDepth !== "undefined") && (a += " " + window.screen.colorDepth), typeof(window.screen.pixelDepth !== "undefined") && (a += " " +
                                        window.screen.pixelDepth));
                                    return a
                                }

                                function nsqddqbd() {
                                    return (new Date).getTimezoneOffset()
                                }

                                function nsbbpd() {
                                    var a = [],
                                        c = /([0-9]+)\.[0-9|.]+/g;
                                    if (window.ActiveXObject) {
                                        if (document.plugins && document.plugins.length > 0)
                                            for (var b = 0; b < document.plugins.length; b++) a.push(document.plugins[b].src.replace(c, "$1"))
                                    } else if (navigator.plugins && navigator.plugins.length > 0)
                                        for (b = 0; b < navigator.plugins.length; b++) a.push(navigator.plugins[b].name.replace(c, "$1"));
                                    a.length > 0 && a.sort();
                                    c = "p";
                                    for (b = 0; b < a.length; b++) c += "," + a[b];
                                    return c
                                }

                                function nspqdqq(a) {
                                    var c = 0,
                                        b = 0,
                                        e, d;
                                    if (a.length === 0) return "00000000";
                                    for (e = 0, l = a.length; e < l; e++) d = a.charCodeAt(e), e % 2 === 0 ? (c = (c << 5) - c + d, c |= 0) : (b = (b << 5) - b + d, b |= 0);
                                    c < 0 && (c = 4294967295 + c + 1);
                                    b < 0 && (b = 4294967295 + b + 1);
                                    return c.toString(16) + b.toString(16)
                                };
                            </script>
                            <script type="text/javascript">
                                function ndpd_load() {
                                    try {
                                        ndwti({
                                            "did": "ndwc",
                                            "fff": "ndpd-%NAME%",
                                            "ffft": "%NAME%",
                                            "ffmm": "fmpm",
                                            "fsss": "ndpd-spbd",
                                            "ddkv": "~~~",
                                            "dddf": "|||",
                                            "ppns": "sssc",
                                            "ppmm": "ensc",
                                            "ppdd": ";",
                                            "ppds": ";",
                                            "wwwe": true,
                                            "fd": {
                                                "s": "1.w-855182.1.2.9YKvTKgr-FK3TwTyftn9mg,,.LlmSxyH6XAw2AVmXrDCEOqIVSmtw1U2ySplf_0hZdYCxh3AGfPw7M74QjWTX72R0VsHsgBdAZ2qkqMNnz4wXnQUFcF9Th7uLNcStLJ5P9EuwLHHU-JfQPczwb_lQ2TXXPEVHGwk2Ig7G54y-lbNULoRuA6UNJYHu3F4Qb76_y0ccRumBQ2SDa4wUPx7na-avDLpb1ulCM9gGHaNK4Z8fET8e3rxqBLRLGAOofWT31tG1Y-wQh095OsNE_oUGIxFSHnNCjIYWqNdQsPYLwXnsz8byGdebcNYOf64fBOJ54GyjsyVY_VBMqyMZgbqklSne9JFmOBX8d-M6luqtJ3GpTA,,",
                                                "f": "1.w-855182.1.2.UP88_ZaW2M7-jGb6WgGAxA,,.HpTmA5peyf8jXpO3Uq1aDEyF5SSS1_TTIqb1nSWtDX3v_R9Dh1bcH7ynXK7qz1S4WsONutyV6oF2Z26SQxkl_fsTc76IVhEh9IrXGOVCLpTASujruPcf7AVqECzfdqjjjyYQHYR4VskpH6BXRxCQQka9yIFf3uJnyNE4dqDVoXH2dcl7GLvRJHzcubAFtYEaCejTU3gY9-FgKNutn8OwYfq2w0EP-eCUFYPCQ7gN_jGBzClaKLszbTlukt5DX947cb_mKsLtgJFAO0329nCzUopCBSXZ5IHy0XtUoPuRrLF2fyqU-FBvFGVB443c7vwlvQO2IQ8oO42Gq7sR5G2TGQ,,",
                                                "fm": "1.w-855182.1.2.QYPj7kTXSD2axDmj9ANeLQ,,.zcuAC5Es7q9j6bobQqaK_w_1lC9vzsKeU_NLxFSN9Dd9p1N1R65rTWo2dpIEUlfoD2QoNx-ubS0qRQK06RUhG0gNxM95yxNUvYSsLK0wC11yRI2Qbo7tRcWdBwV-awmhKuXEg1hCspFiEgZtaYUrkvc1jTSp-eJioT0WCm2SAt08Lt4pTvaqND25jsHyV3kQXzxOkbUp7rBbF4nK6xnB3kyBKqtAgPeWg-dab0YOPhlYs5D0QCHm0m53HD0oMGyNrs_Sm62yFjFHQvTZtLKLfT0hV3quvF5k2Ux18rUA20d_s-wR5-Mha0AP0ixtzGE4RruP6djSsOQINGVMjCn9eA,,",
                                                "w": "1.w-855182.1.2.LkgNzzlVKJj-8fk9PELq7A,,.8ZyOc5TUy9N4DLcqRwkVWiTOJQ4RaJp_6AwvYdwtYCuQyM6yBYW8A82nWPoyRGrGn5UaFcnf6GPgFIRaLrF8S6Zzwp1SQJB4z7T7l3Vi5NsF3NV1JgDmaleIxe4v38PmlRp4hfK_j25NSLNov8qerfud3S0FT-2SoOm-L7WIQSDyoAaeWb_4K4F6AxwDDAtsBMOgPULAnYstsC2Iwcw_E4UVRKcnOHscAaux6ZR8ybfGHbUMDbgcKrcYRWk0eoOAGSn4bM2z-O40xcbL-aaEsX9BKlmU8IvIRhKfE_NEbyQ6DaBDJzw1scAy9XbTYs8oRxpzmpcIEFqNaeWWIPu6ca-gCBpCioDmfUckwxK7pWwvK-2Op3mpwxj6Lg1Ws6iz"
                                            },
                                            "wmd": {
                                                "ipr": {
                                                    "fm": ["@"],
                                                    "lm": true
                                                },
                                                "di": [],
                                                "wk": {
                                                    "r": "https:\/\/s.mzstatic.com\/captcha\/api-us-west-2.ndsopaapl.nudatasecurity.com\/1.0\/w\/34897\/w-855182\/w",
                                                    "w": "1.w-855182.1.2.5FaPryJhZ7vJTwtvcLlDBw,,.seWBjwzAW1CT2RRnlW8bpggXbtdn6jTDU6OGc5l41eCLzKgh3oGjaG5tXqO9HYSfoD79lX1s40KX-MAs6tRoXjcqKY0dkzttaaOgVvlnPX6y5Up_zEbXS67VUvYIxGFMaFIFAllh5rrpfUAnOqv1ob0MhIVGlehe331N_eu-78160MqMjcn6vYq9dA-NiUV2dPpaV0vB27n9t2ch8Nj6I5VNGK8bxXTx0qhqY8BrvZAdmkuAic5XPFQfm7vfeoS10dposx7H6_P2nMQkJHNl-9i1oRiYQwFKeCUIlAqoHWxke_Xum_gdkR1uDIzRU3S-7K_rV12qPxz8_7eli3HapFz3JlhnTWRbpsluqpyAaUk,"
                                                }
                                            }
                                        });
                                    } catch (err) {
                                        var ndpd_suffixes = [""];
                                        for (var ndpd_i = ndpd_suffixes.length - 1; ndpd_i >= 0; --ndpd_i) {
                                            var ndpd_suffix = ndpd_suffixes[ndpd_i];
                                            var ndpd_jse = document.createElement("input");
                                            ndpd_jse.type = "hidden";
                                            ndpd_jse.name = "ndpd-jse";
                                            ndpd_jse.value = err.message;
                                            ndpd_jse.id = "ndpd-jse" + ndpd_suffix;
                                            document.getElementById("ndwc" + ndpd_suffix).appendChild(ndpd_jse);
                                        }
                                    }
                                }
                                if ("complete" === document.readyState) {
                                    ndpd_load();
                                } else if (window.addEventListener) {
                                    window.addEventListener("load", ndpd_load);
                                } else if (window.attachEvent) {
                                    window.attachEvent("onload", ndpd_load);
                                } else {
                                    ndpd_load();
                                }
                            </script>
                            <input aria-hidden="true" id="hidden-captcha-player-mode" type="hidden" value="VIDEO" name="captchaMode" />
                        </div>



                    </div>
                    <div class="account-errors">
                        <ul class="page-error" role="presentation">





                        </ul>
                    </div>
                </div>

                <div class="footer-instruction">

                    <p>Apple 使用符合行业标准的加密方法保护您个人信息的机密性。</p>

                </div>

            </section>





            <nav role="presentation">
                <input class="back" type="submit" value="返回" name="2.0.1.1.3.0.7.11.3.3" />
                <input class="cancel" type="submit" value="取消" name="2.0.1.1.3.0.7.11.3.5" />


                <input class="continue" title="创建 Apple ID" type="submit" value="创建 Apple ID" name="2.0.1.1.3.0.7.11.3.9.1" />

            </nav>
            <input id="machineGUID" class="optional" name="machineGUID" type="hidden" value="" />
            <input id="pageUUID" class="optional" name="mzPageUUID" type="hidden" value='giW+01s3kw/0CXkPxat/OUHW1A0=' />
            <input id="signature" class="optional" name="xAppleActionSignature" type="hidden" value="" />
            <input id="longName" class="optional" name="longName" type="hidden" value="" />
        </form>



        <div class="simple-footer">
            <div class="legal">
                <span class="copyright">Copyright &copy; 2016 Apple Inc. <a target="_blank" href="https://www.apple.com/cn/legal/">保留所有权利。</a></span>
                <ul metrics-loc="Legal" role="presentation">
                    <li metrics-loc="Link_"><a target="_blank" href="https://www.apple.com/cn/privacy/">隐私政策</a></li>
                    <li metrics-loc="Link_"><a target="_blank" href="https://www.apple.com/legal/itunes/ww/">服务条款</a></li>
                </ul>
            </div>
        </div>





    </div>

</body>


</html>
`