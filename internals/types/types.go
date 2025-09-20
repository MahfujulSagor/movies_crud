package types

type Movie struct {
	ID       int64     `json:"id"`
	Title    string    `json:"name" validate:"required"`
	Rating   int       `json:"rating" validate:"required,gte=0,lte=10"`
	Director *Director `json:"director"`
	Cast     *Cast     `json:"cast"`
}

type Director struct {
	ID   int64  `json:"id"`
	Name string `json:"name" validate:"required"`
	Age  int    `json:"age" validate:"required,gte=0,lte=110"`
}

type Cast struct {
	ID      int64  `json:"id"`
	Actor   string `json:"actor" validate:"required"`
	Actress string `json:"actress" validate:"required"`
}
