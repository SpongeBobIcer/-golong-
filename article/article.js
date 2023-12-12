function loadArticles(page) {
    // 使用 Fetch API 向你的 Go 服务器发起 AJAX 请求
    fetch("http://localhost:8080/showArticles?page=" + page, {
        method: "GET"
    })
    .then(response => response.json())
    .then(data => {
        // 处理接收到的数据
        displayArticles(data);
    })
    .catch(error => {
        console.error("Error fetching articles:", error);
    });
}

function displayArticles(articles) {
    const articlesList = document.querySelector(".articles-list");

    // 清空之前的内容
    articlesList.innerHTML = '';

    // 显示新的文章
    articles.forEach(article => {
        const articleElement = document.createElement("div");
        
        // 创建带有链接的标题元素
        const titleElement = document.createElement("a");
        titleElement.textContent = article.title;
        titleElement.href = `articleContent.html?id=${article.id}`; // 将文章ID作为查询参数传递给文章内容页面
        titleElement.target = "_blank"; // 在新窗口中打开链接
        
        articleElement.appendChild(titleElement);

        articlesList.appendChild(articleElement);
    });
}
function loadPrev(){
    if(currentPage > 1){
        currentPage--;
        loadArticles(currentPage);
    }
}
function loadMore() {
    // 在这里调用 loadArticles 函数来加载更多文章
    // 你需要跟踪当前页数，可能使用一个全局变量或其他方式
    // 示例：每次加载增加一页
    currentPage++;
    loadArticles(currentPage);
}

// 初始化页面时加载第一页的文章
let currentPage = 1;
loadArticles(currentPage);