package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func GetYoutubeMovieInfo() ([]map[string]interface{}, error){
	ctx := context.Background()
	apiKey := getYoutubeKeyFromJson()
	if apiKey == "" {
		return nil, fmt.Errorf("key.jsonからyoutube API keyを取得できませんでした")
	}

	service, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Error creating YouTube client: %v", err)
	}

	channelIdList, err := getChannelIdFromUserName(service)
	if err != nil {
		return nil, err
	}

	infoList := []map[string]interface{}{}

	for _, channelId := range channelIdList {
		call := service.Search.List([]string{"id", "snippet"}).
			ChannelId(channelId).Order("date").
			MaxResults(3)

		response, err := call.Do()
		if err != nil {
			log.Fatalf("Error making search API call: %v", err)
		}

		for _, item := range response.Items {
			info := map[string]interface{}{}
			info["videoId"] = item.Id.VideoId
			info["title"] = item.Snippet.Title
			info["thumbnailURL"] = item.Snippet.Thumbnails.High.Url
			info["publishedAt"] = item.Snippet.PublishedAt
			info["userName"] = item.Snippet.ChannelTitle
			infoList = append(infoList, info)
		}
	}
	return infoList, nil

}

// ユーザ名からチャンネルIDを取得する
func getChannelIdFromUserName(service *youtube.Service) ([]string, error) {

	userNameList, err := getUserNameFromJSON()
	if err != nil {
		return nil, err
	}

	channelIdList := []string{}

	for _, userName := range userNameList {
		call := service.Search.List([]string{"snippet"}).Q(userName).MaxResults(1)
		response, err := call.Do()

		if err != nil {
			log.Fatalf("Error making channel search API call: %v", err)
		}

		for _, item := range response.Items {
			channelIdList = append(channelIdList, item.Id.ChannelId)
		}
	}

	return channelIdList, nil

}

func getUserNameFromJSON() ([]string, error) {
	file, err := os.Open("users.json")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer file.Close()

	users := map[string]interface{}{}

	if err := json.NewDecoder(file).Decode(&users); err != nil {
		fmt.Println(err)
		return nil, err
	}

	usersList, ok := users["youtube_user_name_list"].([]interface{})
	if !ok {
		fmt.Println(usersList)
		fmt.Println("users.jsonがおかしいです")
		return nil, err
	}

	usersListStr := []string{}
	for _, user := range usersList {
		userStr, ok := user.(string)
		if !ok {
			fmt.Println("users.jsonがおかしいです")
			return nil, err
		}
		usersListStr = append(usersListStr, userStr)
	}
	fmt.Println("userNameListStr: ", usersListStr)

	return usersListStr, nil
}

func getYoutubeKeyFromJson() string {
	file, err := os.Open("key.json")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer file.Close()

	key := map[string]interface{}{}

	if err := json.NewDecoder(file).Decode(&key); err != nil {
		fmt.Println(err)
		return ""
	}

	api_key, ok := key["youtube_api_key"].(string)

	if !ok {
		fmt.Println("key.jsonにyoutube API keyがありません")
		return ""
	}

	return api_key
}