package libs

import (
	"os"
	"strings"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

func NewSmsg(mess *events.Message, sock *NewClientImpl) *IMessage {
	var command string
	var media whatsmeow.DownloadableMessage
	var isOwner = false
	var owner []string
	botNum, _ := sock.ParseJID(sock.WA.Store.ID.User)
	quotedMsg := mess.Message.GetExtendedTextMessage().GetContextInfo().GetQuotedMessage()
	owner = append(owner, os.Getenv("Owner_Number"))
	owner = append(owner, botNum.String())

	for _, own := range owner {
		if own == mess.Info.Sender.ToNonAD().String() {
			isOwner = true
		}
	}

	if pe := mess.Message.GetExtendedTextMessage().GetText(); pe != "" {
		command = pe
	} else if pe := mess.Message.GetImageMessage().GetCaption(); pe != "" {
		command = pe
	} else if pe := mess.Message.GetVideoMessage().GetCaption(); pe != "" {
		command = pe
	}

	if quotedMsg != nil && (quotedMsg.ImageMessage != nil || quotedMsg.VideoMessage != nil || quotedMsg.StickerMessage != nil) {
		if msg := quotedMsg.GetImageMessage(); msg != nil {
			media = msg
		} else if msg := quotedMsg.GetVideoMessage(); msg != nil {
			media = msg
		} else if msg := quotedMsg.GetStickerMessage(); msg != nil {
			media = msg
		}
	} else if mess.Message != nil && (mess.Message.ImageMessage != nil || mess.Message.VideoMessage != nil) {
		if msg := mess.Message.GetImageMessage(); msg != nil {
			media = msg
		} else if msg := mess.Message.GetVideoMessage(); msg != nil {
			media = msg
		}
	} else {
		media = nil
	}

	return &IMessage{
		From:        mess.Info.Chat,
		Sender:      mess.Info.Sender,
		PushName:    mess.Info.PushName,
		OwnerNumber: owner,
		IsOwner:     isOwner,
		IsBot:       mess.Info.IsFromMe,
		IsGroup:     mess.Info.IsGroup,
		Querry:      strings.Join(strings.Split(command, " ")[1:], ` `),
		Body:        command,
		Command:     strings.ToLower(strings.Split(command, " ")[0]),
		Media:       media,
		IsImage: func() bool {
			if mess.Message.GetImageMessage() != nil {
				return true
			} else {
				return false
			}
		}(),
		IsAdmin: func() bool {
			if !mess.Info.IsGroup {
				return false
			}
			admin, err := sock.FetchGroupAdmin(mess.Info.Chat)
			if err != nil {
				return false
			}
			for _, v := range admin {
				if v == mess.Info.Sender.String() {
					return true
				}
			}
			return false
		}(),
		QuotedMsg: mess.Message.GetExtendedTextMessage().GetContextInfo(),
		ID: &waProto.ContextInfo{
			StanzaId:      &mess.Info.ID,
			Participant:   proto.String(mess.Info.Sender.String()),
			QuotedMessage: mess.Message,
		},
		IsQuotedImage: func() bool {
			if quotedMsg.GetImageMessage() != nil {
				return true
			} else {
				return false
			}
		}(),
		IsQuotedSticker: func() bool {
			if quotedMsg.GetStickerMessage() != nil {
				return true
			} else {
				return false
			}
		}(),
		Reply: func(text string) {
			sock.SendText(mess.Info.Chat, text, &waProto.ContextInfo{
				StanzaId:      &mess.Info.ID,
				Participant:   proto.String(mess.Info.Sender.String()),
				QuotedMessage: mess.Message,
			})
		},
	}
}
