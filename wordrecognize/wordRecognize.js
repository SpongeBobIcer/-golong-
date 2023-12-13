var data = {}; // 创建一个空对象来存储数据
// 更新表格数据
function updateTable(data, tableId, content) {
  const table = document.getElementById(tableId);
  // 清空表格内容
  table.innerHTML = "";
  // 提示文本
  const tipText = document.getElementById("tip");
  tipText.innerHTML = "";

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
      row.id = `word-row-${rowData.wordID}`; // 为行设置唯一ID
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
function getNewWordAndTranslation() {
  var request = new XMLHttpRequest();
  request.open("GET", "http://localhost:8080/getRandomWord");

  request.onreadystatechange = function () {
    if (request.readyState === 4 && request.status === 200) {
      // 获取到新的单词和翻译
      var response = JSON.parse(request.responseText);

      // 更新页面上的单词和翻译
      document.getElementById("word").textContent = response.word;
      document.getElementById("translation").textContent = response.translation;

      // 更新数据
      data = response;
    }
  };

  request.send();
}

function addToEasyWord(wordID) {
  const requestBody = JSON.stringify({ wordID });
  var userToken = localStorage.getItem("token");
  fetch("http://localhost:8080/addToEasyWord", {
    method: "POST",
    body: requestBody,
    headers: {
      "Content-Type": "application/json",
      Authorization: "Bearer " + userToken,
    },
  })
    .then((response) => {
      if (response.ok) {
        console.log("Word added to easy words!");
      } else {
        console.error("Failed to add word to easy words");
      }
    })
    .catch((error) => {
      console.error("Error adding word to easy words:", error);
    });
}
function addToErrorWord(wordID) {
  const requestBody = JSON.stringify({ wordID });
  var userToken = localStorage.getItem("token");
  fetch("http://localhost:8080/addToErrorWord", {
    method: "POST",
    body: requestBody,
    headers: {
      "Content-Type": "application/json",
      Authorization: "Bearer " + userToken,
    },
  })
    .then((response) => {
      if (response.ok) {
        console.log("Word added to error words!");
      } else {
        console.error("Failed to add word to error words");
      }
    })
    .catch((error) => {
      console.error("Error adding word to error words:", error);
    });
}
// 在页面加载时执行此代码块
window.addEventListener("load", function () {
  fetch("http://localhost:8080/getRandomWord")
    .then((response) => response.json())
    .then((initialData) => {
      data = initialData;
      document.getElementById("word").textContent = data.word;
      const translationElement = document.getElementById("translation");
      const wordElement = document.getElementById("word");
      translationElement.textContent = data.translation;
      translationElement.style.display = "none";
      const spellInput = document.getElementById("spelling-input");
      spellInput.style.display = "none";
      const submitButton = document.querySelector(".submit-button");
      submitButton.style.display = "none";
      const simpleWordButton = document.querySelector(".simple-word-button");
      simpleWordButton.style.display = "inline";
      // 当“认识”按钮被点击时
      document
        .querySelector(".know-button")
        .addEventListener("click", function () {
          translationElement.style.display = "inline";
          wordElement.style.display = "none";
          spellInput.style.display = "inline";
          submitButton.style.display = "inline";
          simpleWordButton.style.display = "inline";
        });

      // 当“不认识”按钮被点击时
      document
        .querySelector(".dont-know-button")
        .addEventListener("click", function () {
          translationElement.style.display = "inline";
          wordElement.style.display = "inline";
          simpleWordButton.style.display = "none";
        });

      // 当“简单词”按钮被点击时
      document
        .querySelector(".simple-word-button")
        .addEventListener("click", function () {
          translationElement.style.display = "inline";
          addToEasyWord(data.wordID);
          layMsg("已加入简单词，不会再出现(可以在简单词界面移除)", { icon: 1 });
          getNewWordAndTranslation();
          translationElement.style.display = "none";
          wordElement.style.display = "inline";
          spellInput.style.display = "none";
          submitButton.style.display = "none";
        });

      // 提交按钮的逻辑
      submitButton.addEventListener("click", function () {
        const userInput = spellInput.value;
        if (userInput === wordElement.textContent) {
          layMsg("拼写正确！", { icon: 1 });
          getNewWordAndTranslation();
          wordElement.style.display = "none";
          spellInput.style.display = "none";
          submitButton.style.display = "none";
        } else {
          // 错误词
          layMsg(
            "拼写错误。正确拼写是: " + wordElement.textContent,
            { icon: 2, time: 1500 },
            () => {
              addToErrorWord(data.wordID);
              layMsg("已加入错误词(可以在错误词界面查看)", {
                icon: 2,
                time: 1500,
              });
            }
          );
          getNewWordAndTranslation();
        }
        translationElement.style.display = "none";
        wordElement.style.display = "inline";
        spellInput.style.display = "none";
        submitButton.style.display = "none";
        simpleWordButton.style.display = "inline";
      });
    })
    .catch((error) => {
      console.error("Error fetching random word:", error);
    });
});
