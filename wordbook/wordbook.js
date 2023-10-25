var userToken = localStorage.getItem("token");
// 创建删除按钮点击事件处理程序
function deleteWordsByType(content, wordID) {
    // 向服务器发送请求以删除单词
    // 如果删除成功，则从表格中删除行

    fetch(`http://localhost:8080/deleteWord?content=${content}&wordID=${wordID}`, {
        method: "DELETE",
        headers: {
            "Authorization": "Bearer " + userToken,
            "Content-Type": "application/json"
        }
    })
    .then(response => {
        if (response.status === 200) {
            // 从表格中删除行
            const table = document.getElementById("word-table");
            const row = document.getElementById(`word-row-${wordID}`);
            table.deleteRow(row.rowIndex);
        } else {
            console.error("Error deleting word:", response.statusText);
        }
    })
    .catch(error => {
        console.error("Error deleting word:", error);
    });
}

// 更新表格数据
function updateTable(data, tableId, content) {
    const table = document.getElementById(tableId);
    // 清空表格内容
    table.innerHTML = "";

    // 创建表头
    const headerRow = table.insertRow(0);
    const wordHeader = headerRow.insertCell(0);
    const translationHeader = headerRow.insertCell(1);
    const errorCountHeader = headerRow.insertCell(2);
    const deleteHeader = headerRow.insertCell(3);
    wordHeader.innerHTML = "单词";
    translationHeader.innerHTML = "翻译";
    errorCountHeader.innerHTML = "错误次数";
    deleteHeader.innerHTML = "操作";

    // 检查 data 是否有效
    if (data && data.length > 0) {
        // 填充数据
        for (let i = 0; i < data.length; i++) {
            const rowData = data[i];
            const row = table.insertRow(i + 1);
            row.id = `word-row-${rowData.wordID}`;  // 为行设置唯一ID
            const wordCell = row.insertCell(0);
            const translationCell = row.insertCell(1);
            const errorCountCell = row.insertCell(2);
            const deleteCell = row.insertCell(3);

            wordCell.innerHTML = rowData.word;
            translationCell.innerHTML = rowData.translation;
            errorCountCell.innerHTML = rowData.errorCount;

            // 创建删除按钮
            const deleteButton = document.createElement("button");
            deleteButton.textContent = "删除";
            deleteButton.classList.add("delete-button");
            deleteButton.addEventListener("click", () => {
                // 调用删除单词函数并传递内容类型和单词ID
                deleteWordsByType(content, rowData.wordID);
            });
            deleteCell.appendChild(deleteButton);
        }
    } else {
        // 如果 data 无效，显示一条消息或采取其他操作
        const messageRow = table.insertRow(1);
        const messageCell = messageRow.insertCell(0);
        messageCell.colSpan = 4;
        messageCell.innerHTML = "没有可用的数据。";
    }
}

// 修改 getWordsByType 函数来接受 content 参数
function getWordsByType(content) {
    // 设置请求头部，包括 Authorization 头部
    var headers = new Headers({
        "Authorization": "Bearer " + userToken,
        "Content-Type": "application/json"
    });

    // 发送请求到后端以获取指定类型的词汇数据
    fetch(`http://localhost:8080/showWordList?content=${content}`, {
        method: "GET",
        headers: headers
    })
    .then(response => response.json())
    .then(data => {
        // 更新表格中的数据
        updateTable(data, "word-table", content);
    })
    .catch(error => {
        console.error(`Error fetching ${content} words:`, error);
    });
}

// 添加事件监听器来触发获取简单词和错误词的数据
document.getElementById("simple-word-button").addEventListener("click", () => getWordsByType("simpleWords"));
document.getElementById("error-word-button").addEventListener("click", () => getWordsByType("errorWords"));
