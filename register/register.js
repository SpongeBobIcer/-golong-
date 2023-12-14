window.onload = function () {
  document
    .getElementById("registerForm")
    .addEventListener("submit", function (event) {
      event.preventDefault();

      var formData = new FormData(event.target);
      var request = new XMLHttpRequest();
      request.open("POST", "http://localhost:8080/register"); // 后端的注册API端点

      request.onreadystatechange = function () {
        if (request.readyState === 4) {
          if (request.status === 200) {
            // 注册成功
            var response = request.responseText;
            layMsg("在注册过程中: " + response);
            window.location.href = "../login/login.html";
          } else {
            // 注册失败
            layMsg("在注册过程中遇到错误: " + request.statusText);
          }
        }
      };

      request.send(formData);
    });
};
