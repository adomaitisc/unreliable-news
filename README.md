# Unreliable Artificial Intelligence News - UAIN

### Welcome to my dumb side project.

![image](https://github.com/adomaitisc/unreliable-news/blob/main/assets/unreliablea-artificial-intelligence-news.png?raw=true)

#### App description

The app collects the most relevant news about Elon Musk and Twitter straight from the Bloomberg's website.
It uses a Paraphrasing API to reqrite the text ( often making it really messy, but that's the fun... i guess)
On every request to the home page, it fetches the brand new article about Musk, and if it is not already on the databse, it creates a new row.

The scraper was very complicated, as it need to collect one script tag from the returned html document, and strip, replace and do most string function to get the actual text. The content was inside a "body" key in a JSON object.

__________

#### Scraper

- Fetching the url with the most relevant news about Elon Musk
- Collects the first one
- Fetches it.
- Collect title and unloaded image, as well as the script tag eith the content.
- The content of the script tag looked somethink like this: - So, it needed to be cleaned.
    ![image](https://github.com/adomaitisc/unreliable-news/blob/main/assets/code-mess.png?raw=true)
    
__________

#### Tech

- The backend is written in GO, with Astro on the frontend.
- The databse is an instance from Planetscale.
- The styling is made with TailwindCSS.

Starting the backend

<pre><code>$ cd backend
$ go mod tidy
$ go run .</code></pre>

Starting the frontend:

<pre><code>$ cd frontend
$ npm init
$ npm run dev</code></pre>

__________

#### API

- The API used is from [API Layer - Paraphraser API]("https://apilayer.com/")
- It is not reliable for this project, causing weird words to appear in the middle of the text, and also random characters.
