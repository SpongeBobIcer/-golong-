var data = {}; // 创建一个空对象来存储数据

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
            "Authorization": "Bearer " + userToken
        }
    })
    .then(response => {
        if (response.ok) {
            console.log("Word added to easy words!");
        } else {
            console.error("Failed to add word to easy words");
        }
    })
    .catch(error => {
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
            "Authorization": "Bearer " + userToken
        }
    })
    .then(response => {
        if (response.ok) {
            console.log("Word added to error words!");
        } else {
            console.error("Failed to add word to error words");
        }
    })
    .catch(error => {
        console.error("Error adding word to error words:", error);
    });
}
// 在页面加载时执行此代码块
window.addEventListener("load", function () {
    fetch("http://localhost:8080/getRandomWord")
    .then(response => response.json())
    .then(initialData => {
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
        document.querySelector(".know-button").addEventListener("click", function () {
            translationElement.style.display = "inline";
            wordElement.style.display = "none";
            spellInput.style.display = "inline";
            submitButton.style.display = "inline";
            simpleWordButton.style.display = "inline";
        });

        // 当“不认识”按钮被点击时
        document.querySelector(".dont-know-button").addEventListener("click", function () {
            translationElement.style.display = "inline";
            wordElement.style.display = "inline";
            simpleWordButton.style.display = "none";
        });

        // 当“简单词”按钮被点击时
        document.querySelector(".simple-word-button").addEventListener("click", function () {
            translationElement.style.display = "inline";
            addToEasyWord(data.wordID);
            alert("已加入简单词，不会再出现(可以在简单词界面移除)");
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
                alert("拼写正确！");
                getNewWordAndTranslation();
                wordElement.style.display = "none";
                spellInput.style.display = "none";
                submitButton.style.display = "none";
            } else {
                alert("拼写错误。正确拼写是: " + wordElement.textContent);
                //错误词
                addToErrorWord(data.wordID)
                alert("已加入错误词(可以在错误词界面查看)");
                getNewWordAndTranslation();
            }
            translationElement.style.display = "none";
            wordElement.style.display = "inline";
            spellInput.style.display = "none";
            submitButton.style.display = "none";
            simpleWordButton.style.display = "inline";
        });
    })
    .catch(error => {
        console.error("Error fetching random word:", error);
    });
});
