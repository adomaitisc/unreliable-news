# Unreliable Artificial Intelligence News - UAIN

![image]("https://github.com/adomaitisc/unreliable-news/blob/main/assets/unreliablea-artificial-intelligence-news.png?raw=true")

### Welcome to my dumb side project.

The app collects the most relevant news about Elon Musk and Twitter straight from the Bloomberg's website.

It uses a Paraphrasing API to reqrite the text ( often making it really messy, but that's the fun... i guess)

On every request to the home page, it fetches the brand new article about Musk, and if it is not already on the databse, it creates a new row.

#### Tech

The backend is written in GO, with Astro on the frontend.

The databse is an instance from Planetscale.

The styling is made with TailwindCSS.

To run the backend:
`go mod tidy`
`cd backend`
`go run .`

To run the frontend:
`cd frontend`
`npm init`
`npm run dev`

#### API

The API used is from [API Layer - Paraphraser API]("https://apilayer.com/")

It is not reliable for this project, causing weird words to appear in the middle of the text, and also random characters.
