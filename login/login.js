document.getElementById("loginForm").addEventListener("submit", function (event) {
    event.preventDefault();
    var formData = new FormData(event.target);
    var request = new XMLHttpRequest();
    request.open("POST", "http://localhost:8080/login"); // 后端的登录API端点

    request.onreadystatechange = function () {
        if (request.readyState === 4) {
            if (request.status === 200) {
                // 登录成功
                 // 从JSON获取token
                var response = JSON.parse(request.responseText); // 解析JSON字符串
                var token = response.token; // 获取 "token" 键的值
                // 存储用户令牌
                localStorage.setItem("token", token);

                // 检查是否已经存储了令牌
                if (localStorage.getItem("token")) {
                    // 从localStorage中获取令牌
                    var userToken = localStorage.getItem("token");
                    // 输出未解码的令牌
                    console.log("User Token:", userToken);
                
                    // 其他字段...
                } else {
                    console.log("Token not found.");
                }
                

                // 重定向到 home 页面
                window.location.href = "../home/home.html";
            } else if (request.status === 401) {
                // 未授权
                alert("Login failed: Invalid email or password");
            } else {
                // 其他错误
                alert("Login failed: " + request.statusText);
            }
        }
    };

    request.send(formData);
});
