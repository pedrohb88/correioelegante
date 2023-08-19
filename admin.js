let phrases = [
    "Viva São João!",
    "Olha a cobra!",
    "Balancê, balancê!",
    "Dois pra lá, dois pra cá!",
    "Anarriê!",
    "Olha o passo do noivo!",
    "Olha o passo da noiva!",
    "Cavalheiros à direita!",
    "Damas à esquerda!",
    "Vai começar o casamento!",
    "Vem aí o padre!",
    "Cai no laço!",
    "Vai e volta!",
    "Anarriê, minha gente!",
    "Olha a chuva!",
    "Vai o noivo, volta a noiva!",
    "Vai o cortejo!",
    "Vai o noivo buscar a noiva!",
    "Vem aí o delegado!",
    "Olha a prisão!",
    "Dama de honra na roça!",
    "Olha a tesoura!",
    "Olha o buquê!",
    "Olha o lenço!",
    "Vai a saia, volta o lenço!",
    "Vai a saia e o lenço!",
    "Olha a serpentina!",
    "Vai o noivo, volta a noiva, e o casamento se inicia!"
  ]

document.addEventListener("DOMContentLoaded", function() {
    // Function to play text as audio using the Web Speech API
    function playTextAsAudio(text) {
        const synth = window.speechSynthesis;
        const utterance = new SpeechSynthesisUtterance(text);
        const brazilianVoice = synth.getVoices().find(voice => voice.lang === 'pt-BR');
        utterance.voice = brazilianVoice;
    
        // Adjust pitch and rate for a more natural sound
        utterance.pitch = 1; // 0 to 2, default: 1
        utterance.rate = 1.2; // 0.1 to 10, default: 1
        synth.speak(utterance);
    }

    // Function to fetch messages and play them as audio
    function fetchAndPlayMessages() {
        fetch("/messages?lastReadMsgID=0") // Adjust the URL as needed
            .then(response => response.json())
            .then(data => {
                if (data.length > 0) {
                    data.forEach(message => {
                        const fullMessage = `Mensagem de ${message.from} para ${message.to}: ${message.text}`;
                        playTextAsAudio(fullMessage);
                        document.getElementById("messageList").innerHTML = document.getElementById("messageList").innerHTML + "<li>" +fullMessage+"</li>"
                    });
                } else {
                    playTextAsAudio(phrases[Math.floor(Math.random() * phrases.length)])
                }
            })
            .catch(error => {
                console.error("Error fetching messages:", error);
            });
    }

    // Fetch and play messages every 10 minutes
    setInterval(fetchAndPlayMessages, 10 * 60 * 1000); // 10 minutes in milliseconds
});
