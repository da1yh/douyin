package controller

var DemoVideoList = []Video{
	{
		Id:            1,
		Author:        DemoUser,
		PlayUrl:       "https://www.w3schools.com/html/movie.mp4",
		CoverUrl:      "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg",
		FavoriteCount: 20,
		CommentCount:  10,
		IsFavorite:    true,
		Title:         "bear",
	},
}

var DemoUser = User{
	Id:            1,
	Name:          "ZywOo",
	FollowCount:   100,
	FollowerCount: 200,
	IsFollow:      true,
}
