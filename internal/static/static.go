package static

var Template = `
<style>
body {
	font-family: Arial, sans-serif;
	background-color: #f4f4f4;
	margin: 0;
	padding: 0;
	display: flex;
	justify-content: center;
	align-items: center;
	min-height: 100vh;
}
.container {
	text-align: center;
	padding: 20px;
	background-color: #ffffff;
	box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
	border-radius: 8px;
	max-width: 600px; 
}
h1 {
	color: #333333;
	margin-bottom: 10px;
}
p {
	color: #666666;
	margin-bottom: 20px;
	word-wrap: break-word; 
	font-size: 18px;
	line-height: 1.5;
}
img {
	max-width: 100%;
	max-height: 600px; /* Adjust this value as needed */
	border-radius: 4px;
	box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}
</style>

<body>
    <div class="container">
        <h1>Random Text and Image</h1>
        <p id="typed-text"></p>
        <img src="data:image/jpeg;base64,{{ .Image }}" alt="Random Image">
    </div>
    <script>
        const textElement = document.getElementById("typed-text");
        const textToType = ` + "`{{.Text}}`" + `;
        
        function typeText(text, element) {
            element.textContent = ""; // Clear existing text
            let currentIndex = 0;

            function typeNextLetter() {
                if (currentIndex < text.length) {
                    element.textContent += text[currentIndex];
                    currentIndex++;
                    setTimeout(typeNextLetter, 30); // Adjust typing speed here
                }
            }

            typeNextLetter();
        }

        typeText(textToType, textElement);
    </script>
</body>
</html>
`
