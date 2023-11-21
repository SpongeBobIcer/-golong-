// articleContent.js

// 获取文章ID
const urlParams = new URLSearchParams(window.location.search);
const articleID = urlParams.get("id");

// 获取文章内容
function loadArticleContent() {
    fetch(`http://localhost:8080/showArticleContent?id=${articleID}`)
        .then(response => response.json())
        .then(data => {
            // 渲染文章内容
            displayArticleContent(data);
        })
        .catch(error => {
            console.error("Error fetching article content:", error);
        });
}

function displayArticleContent(article) {
    // 显示文章标题
    document.getElementById("articleTitle").textContent = article.title;

    // 显示文章内容
    document.getElementById("articleContent").innerHTML = article.content;
}

// 初始化页面时加载文章内容
loadArticleContent();
