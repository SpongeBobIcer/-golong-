// 获取用户数据
function getUserData() {
    // 从 localStorage 中获取用户令牌
    var userToken = localStorage.getItem("token");

    // 设置请求头部，包括 Authorization 头部
    var headers = new Headers({
        "Authorization": "Bearer " + userToken,
        "Content-Type": "application/json"
    });

    // 发送请求到后端以获取用户数据
    fetch("http://localhost:8080/showUserData", {
        method: "GET",
        headers: headers
    })
    .then(response => response.json())
    .then(data => {
        // 更新页面上的用户信息
        console.log("user data:", data);
        document.getElementById("userID").textContent = data.userInfo.UserID;
        document.getElementById("userName").textContent = data.userInfo.Username;
        document.getElementById("userEmail").textContent = data.userInfo.Email;
        document.getElementById("userPhone").textContent = data.userInfo.PhoneNumber;
        document.getElementById("userLevel").textContent = data.userData.Level
        document.getElementById("totalProgess").textContent = data.userData.TotalProgess
        document.getElementById("dailyProgess").textContent = data.userData.DailyProgess
        document.getElementById("totalGoal").textContent = data.userData.TotalGoal
        document.getElementById("dailyGoal").textContent = data.userData.DailyGoal
        document.getElementById("accuracy").textContent = data.userData.Accuracy
        document.getElementById("learningCount").textContent = data.userData.LearningCount
    })
    .catch(error => {
        console.error("Error fetching user data:", error);
    });
}

// 添加事件监听器，当页面加载时获取用户数据
window.addEventListener("load", () => {
    getUserData();
});
document.addEventListener("DOMContentLoaded", function() {
    var userToken = localStorage.getItem("token");
    const changePasswordButton = document.getElementById("changePasswordButton");

    changePasswordButton.addEventListener("click", function() {
        // 创建原密码输入框
        const oldPasswordInput = document.createElement("input");
        oldPasswordInput.type = "password";
        oldPasswordInput.placeholder = "原密码";
        oldPasswordInput.id = "oldPasswordInput";
        
        // 创建新密码输入框
        const newPasswordInput = document.createElement("input");
        newPasswordInput.type = "password";
        newPasswordInput.placeholder = "新密码";
        newPasswordInput.id = "newPasswordInput";

        // 创建确认新密码输入框
        const confirmNewPasswordInput = document.createElement("input");
        confirmNewPasswordInput.type = "password";
        confirmNewPasswordInput.placeholder = "确认新密码";
        confirmNewPasswordInput.id = "confirmNewPasswordInput";

        // 创建提交按钮
        const submitButton = document.createElement("button");
        submitButton.textContent = "提交";
        submitButton.id = "submitButton";

        // 添加输入框和提交按钮到用户信息部分
        changePasswordButton.style.display = "none"; // 隐藏"修改密码"按钮
        document.querySelector(".user-info").appendChild(oldPasswordInput);
        document.querySelector(".user-info").appendChild(newPasswordInput);
        document.querySelector(".user-info").appendChild(confirmNewPasswordInput);
        document.querySelector(".user-info").appendChild(submitButton);

        // 监听提交按钮的点击事件
        submitButton.addEventListener("click", function() {
            // 获取原密码、新密码和确认新密码输入框的值
            const oldPassword = document.getElementById("oldPasswordInput").value;
            const newPassword = document.getElementById("newPasswordInput").value;
            const confirmNewPassword = document.getElementById("confirmNewPasswordInput").value;

            // 构建请求数据
            const passwordData = {
                oldPassword: oldPassword,
                newPassword: newPassword,
                confirmNewPassword: confirmNewPassword
            };

            // 发送POST请求到后端，将密码数据放入请求正文
            fetch("http://localhost:8080/changePassword", {
                method: "POST",
                headers: {
                    "Authorization": "Bearer " + userToken,
                    "Content-Type": "application/json"
                },
                body: JSON.stringify(passwordData) // 将密码数据转换为 JSON 字符串
            })
            .then(response => {
                if (response.status === 200) {
                    // 密码修改成功
                    layMsg("密码修改成功", {icon: 1});
                } else {
                    layMsg("密码修改失败", {icon: 2});
                }
            })
            .catch(error => {
                console.error("Error changing password:", error);
            });
        });

        // 移除"修改密码"按钮
        changePasswordButton.style.display = "none";
    });
});
