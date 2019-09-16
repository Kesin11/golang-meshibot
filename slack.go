package main

import (
	"fmt"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/nlopes/slack"
)

const botImageURL = "https://pbs.twimg.com/profile_images/647226377066840064/JPt0F01G_200x200.png"

func buildRetaurantBlock(r Restaurant) *slack.SectionBlock {
	infoText := fmt.Sprintf("*<%v|%v>*\n%v\n*<%v|ランチメニュー>*", r.URL, r.Name, r.Description, r.LunchURL)
	info := slack.NewTextBlockObject("mrkdwn", infoText, false, false)
	image := slack.NewImageBlockElement(r.ImageURL, "thumbnail")

	section := slack.NewSectionBlock(info, nil, slack.NewAccessory(image))
	return section
}

func buildMsgOptionBlock(restaurants []Restaurant) slack.MsgOption {
	// Build Header Section Block, includes text and overflow menu
	headerText := slack.NewTextBlockObject("mrkdwn", "お店が見つかりました", false, false)

	// Create the header section
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	// Shared Divider
	divSection := slack.NewDividerBlock()

	blocks := []slack.Block{headerSection}
	for _, restaurant := range restaurants {
		blocks = append(blocks, divSection)
		blocks = append(blocks, buildRetaurantBlock(restaurant))
	}

	return slack.MsgOptionBlocks(blocks...)
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
	keyword := strings[1]
	restaurants, err := s.restaurantClient.fetchRandom(keyword, 5)
	if err != nil {
		return fmt.Errorf("failed to fetch restaurant: %s", err)
	}
	if len(restaurants) < 1 {
		rtm.SendMessage(rtm.NewOutgoingMessage("1件も見つかりませんでした", channelID))
		return nil
	}

	msgOptionBlock := buildMsgOptionBlock(restaurants)
	msgOptionIconURL := slack.MsgOptionIconURL(botImageURL)

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

// RestaurantClient レストラン取得のAPIクライアントを差し替え可能にするため
type RestaurantClient interface {
	fetchRandom(keyword string, limit int) ([]Restaurant, error)
}

// SlackListener RTMとbotからの返信を扱う
type SlackListener struct {
	client           *slack.Client
	botUserID        string
	restaurantClient RestaurantClient
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
