function initTheme() {
    const theme = sessionStorage.getItem('theme') || 'light';
    document.body.className = theme;
}

initTheme();

window.addEventListener('storage', () => {
    initTheme();
});

window.addEventListener('load', () => {
    const data = JSON.parse(sessionStorage.getItem("searchData"));
    
    if (data) {
        console.log(data);
        
        let title = document.getElementById("titre");
        title.textContent = data.word;
        
        const section = document.getElementById("links");
        
        if (data.link && Array.isArray(data.link)) {
            data.link.forEach((link, index) => {
                const card = document.createElement("div");
                card.className = 'result-card';
                
                const domain = extractDomain(link);
                const tfIdf = data.tfIdf || 'N/A';
                const world = data.world || '';
                
                const a = document.createElement("a");
                a.href = link;
                a.target = "_blank";
                a.className = 'result-link';
                a.textContent = truncateText(link, 70);
                
                card.innerHTML = `
                    <div class="result-rank">#${index + 1}</div>
                    <div class="result-url">${domain}</div>
                    <div class="result-description">${world || 'Résultat de recherche'}</div>
                    <div class="result-metadata">
                        <span class="result-badge">Score: ${(parseFloat(tfIdf) * 100).toFixed(2)}%</span>
                        <span>Catégorie: ${world || 'Général'}</span>
                    </div>
                `;
                
                const linkWrapper = document.createElement("div");
                linkWrapper.appendChild(a);
                card.insertBefore(linkWrapper, card.childNodes[1]);
                
                section.appendChild(card);
            });
        }
    }
});

function extractDomain(url) {
    try {
        const urlObj = new URL(url);
        return urlObj.hostname;
    } catch {
        return url;
    }
}

function truncateText(text, maxLength) {
    return text.length > maxLength ? text.substring(0, maxLength) + '...' : text;
}

document.getElementById('btn').addEventListener('click', () => {
    const query = document.getElementById('inp').value.trim();
    if (query) {
        const packag = {
            action: "searchGo",
            data: { value: query }
        };

        fetch("http://127.0.0.1:8080/back/main", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(packag)
        })
        .then(response => {
            if (!response.ok)
                throw new Error("HTTP Error " + response.status + " " + response.statusText);
            return response.json();
        })
        .then(data => {
            console.log("Réponse du serveur :", data);
            sessionStorage.setItem("searchData", JSON.stringify(data));
            window.location.href = "inSearch.html";
        })
        .catch(error => {
            console.log(`${query} -> not in bdd`);
        });
    }
});

document.getElementById('inp').addEventListener('keypress', (e) => {
    if (e.key === 'Enter') {
        document.getElementById('btn').click();
    }
});