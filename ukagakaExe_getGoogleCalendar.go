package main

import (
    "fmt"
    "flag"
    "time"
    "os"
    "context"
    "regexp"
    "strconv"

    "google.golang.org/api/calendar/v3"
    "google.golang.org/api/option"
)


//credentials file path
//gmail.com
//day or week
//time zone
//sep
func main()  {
    flag.Parse()
    args := flag.Args()
    if len( args ) != 5 {
        fmt.Println( "引数の数がおかしい。" )
        return
    }
    credentialFilePath  := args[0]
    Gmail               := args[1]
    Target              := args[2]
    //TimeZoneStr            := 3
    TimeZoneStr            := args[3]
    Sep                 := args[4]

    TimeZone ,_  := strconv.Atoi( TimeZoneStr ) 


    //認証ファイルチェック
    _, err := os.Stat( credentialFilePath )
    if err != nil {
        fmt.Println( "認証用ファイルがないよ。" )
        return
    }
    ctx := context.Background()
    calendarService, err := calendar.NewService(ctx, option.WithCredentialsFile( credentialFilePath ))
    if err != nil { 
        fmt.Println( "jsonfile load error" ) 
        return 
    }


    var MIN string = ""
    var MAX string = ""
    today       := time.Date(  time.Now().Year() , time.Now().Month() , time.Now().Day() , TimeZone , 0 , 0 , 0 , time.UTC )
    tomorrow    := today.AddDate( 0, 0, 1)
    week        := tomorrow.AddDate( 0, 0, 7)

    if Target == "day" {
        MIN = today.Format(time.RFC3339)
        MAX = tomorrow.Format(time.RFC3339)
    } else if Target == "week" {
        MIN = tomorrow.Format(time.RFC3339)
        MAX = week.Format(time.RFC3339)
    } else {
        fmt.Println( "Target :" + Target + ". Target is day or week" )
        return
    }


    events, err := calendarService.Events.
        List(           Gmail       ).
        TimeMin(        MIN         ).
        TimeMax(        MAX         ).
        OrderBy(        "startTime" ).
        SingleEvents(   true        ).
        Do()
    if err != nil {
        fmt.Printf( "check email : %v" , err ) 
        return
    }


    //予定が見つからなかった場合。
    if len( events.Items ) == 0 {
        //fmt.Println( "予定なし。" )
        return
    }
    res := ""
    for _,item := range events.Items {

        title := item.Summary
        if title == "" {
            title = "No Title"
        }

        //終日
        timeText := ""
        startTime := item.Start.Date 
        if startTime != "" {
            x := regexp.MustCompile( "^........" )
            startTime = x.ReplaceAllString( startTime , "" )
            timeText = startTime + "日 : " + title

        //指定時間あり。
        } else {
            startTime = item.Start.DateTime
            x := regexp.MustCompile( "^........")
            day := x.ReplaceAllString( startTime , "" )
            x = regexp.MustCompile( "^..")
            dayx := x.FindAllStringSubmatch( day , 1 )
            //fmt.Println( day )
            //fmt.Println( dayx )

            x = regexp.MustCompile( "^..........." )
            hhmm := x.ReplaceAllString( startTime , "" )
            x = regexp.MustCompile( `:00\+.*?$` )
            hhmm = x.ReplaceAllString( hhmm , "" )
            //timeText = "日 " + hhmm + " : " + title
            timeText = dayx[0][0] + "日 " + hhmm + " : " + title
        }
        res = res + timeText + Sep
    }
    fmt.Println( res )
}





















