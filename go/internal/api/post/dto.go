package post

type getPostResponse struct {
	Title     string `json:"Title"`
	SK        string `json:"SK"`
	Content   string `json:"Content"`
	CreatedAt string `json:"CreatedAt"`
}

func toGetPostResponse(p post) getPostResponse {
	return getPostResponse{
		Title:     p.Title,
		Content:   p.Content,
		SK:        p.SK,
		CreatedAt: p.CreatedAt,
	}
}

func toGetPostsResponse(posts []post) []getPostResponse {
	res := make([]getPostResponse, len(posts))

	for i, p := range posts {
		res[i] = getPostResponse{
			Title:     p.Title,
			Content:   p.Content,
			SK:        p.SK,
			CreatedAt: p.CreatedAt,
		}
	}

	return res
}
