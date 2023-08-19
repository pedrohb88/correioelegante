document.getElementById("sendButton").addEventListener("click", function() {
    var from = document.getElementById("from").value;
    var to = document.getElementById("to").value;
    var text = document.getElementById("text").value;

    var xhr = new XMLHttpRequest();
    xhr.open("POST", "/messages", true);
    xhr.setRequestHeader("Content-Type", "application/json");

    xhr.onreadystatechange = function() { 
        if (xhr.readyState === XMLHttpRequest.DONE) {
            if (xhr.status === 200) {
                document.getElementById("successMessage").textContent = "Mensagem enviada com sucesso!";
                document.getElementById("from").value = "";
                document.getElementById("to").value = "";
                document.getElementById("text").value = "";
            } else {
                document.getElementById("successMessage").textContent = "Erro ao enviar mensagem.";
            }
        }
    };

    var data = JSON.stringify({ from: from, to: to, text: text });
    xhr.send(data);
});


