package main

import (
	"fmt"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/nlopes/slack"
)

type Hotel struct {
	Name       string
	URL        string
	Rate       string
	RateString string
	Price      string
	ImageURL   string
	Location   string
}

func createHotelBlock(hotel Hotel) (*slack.SectionBlock, *slack.ContextBlock) {
	infoText := fmt.Sprintf("*<%v|%v>*\n%v\n%v\n%v", hotel.URL, hotel.Name, hotel.Rate, hotel.RateString, hotel.Price)
	info := slack.NewTextBlockObject("mrkdwn", infoText, false, false)
	image := slack.NewImageBlockElement(hotel.ImageURL, "thumbnail")
	loc := slack.NewTextBlockObject("plain_text", hotel.Location, true, false)

	section := slack.NewSectionBlock(info, nil, slack.NewAccessory(image))
	context := slack.NewContextBlock("", []slack.MixedElement{loc}...)

	return section, context
}

func buildMsgOptionBlock() slack.MsgOption {
	// Build Header Section Block, includes text and overflow menu
	headerText := slack.NewTextBlockObject("mrkdwn", "お店が見つかりました", false, false)

	// Create the header section
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	// Shared Divider
	divSection := slack.NewDividerBlock()

	// First Hotel Listing
	var hotel = Hotel{
		Name:       "Windsor Court Hotel",
		URL:        "fakeLink.toHotelPage.com",
		Rate:       "★★★★★",
		RateString: "Rated: 9.4 - Excellent",
		Price:      "$340 per night",
		ImageURL:   "https://api.slack.com/img/blocks/bkb_template_images/tripAgent_1.png",
		Location:   "Location: Central Business District",
	}
	hotelOneSection, hotelOneContext := createHotelBlock(hotel)

	// Second Hotel Listing
	hotelTwoSection, hotelTwoContext := createHotelBlock(hotel)

	// Third Hotel Listing
	hotelThreeSection, hotelThreeContext := createHotelBlock(hotel)

	return slack.MsgOptionBlocks(
		headerSection,
		divSection,
		hotelOneSection,
		hotelOneContext,
		divSection,
		hotelTwoSection,
		hotelTwoContext,
		divSection,
		hotelThreeSection,
		hotelThreeContext,
	)
}

func (s *SlackListener) handleMessage(msg slack.Msg, rtm *slack.RTM) error {
	channelID := msg.Channel
	text := msg.Text
	strings := strings.Split(text, " ")

	if !s.isMentionToBot(strings) {
		return nil
	}
	if len(strings) < 2 {
		rtm.SendMessage(rtm.NewOutgoingMessage("検索ワードも入力してください", channelID))
		return nil
	}

	// 検索ワード
	// term := strings[1]

	msgOptionBlock := buildMsgOptionBlock()
	msgOptionIconURL := slack.MsgOptionIconURL("https://pbs.twimg.com/profile_images/1146470842546548737/D9rq59or_200x200.jpg")

	if _, _, err := s.client.PostMessage(channelID, msgOptionBlock, msgOptionIconURL); err != nil {
		return fmt.Errorf("failed to post message: %s", err)
	}

	return nil
}

func (s *SlackListener) isMentionToBot(strings []string) bool {
	if strings[0] == fmt.Sprintf("<@%v>", s.botUserID) {
		return true
	}
	return false
}

// SlackListener RTMとbotからの返信を扱う
type SlackListener struct {
	client    *slack.Client
	botUserID string
}

// ListenAndResponse SlackのRTMを受信して処理を振り分ける
func (s SlackListener) ListenAndResponse() error {
	rtm := s.client.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			spew.Dump(ev)
			s.handleMessage(ev.Msg, rtm)

		case *slack.RTMError:
			return fmt.Errorf("Error: %s", ev.Error())

		case *slack.InvalidAuthEvent:
			return fmt.Errorf("Invalid credentials")

		default:

			// Ignore other events..
			// fmt.Printf("Unexpected: %v\n", msg.Data)
		}
	}

	return nil
}
