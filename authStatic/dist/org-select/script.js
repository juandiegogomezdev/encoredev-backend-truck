

addEventListener("DOMContentLoaded", async (event) => {
    const cardsSpace = document.querySelectorAll(".cardSpace");

    try {
        const response = await fetch(`${window.APP_CONFIG.api_url}/org/session`)
        if (response.ok) {
            const sessionData = await response.json();
            console.log("session data: ", sessionData);
        }
        else {
            console.error("Failed to fetch session data: ", response.status);
        }
    }
    catch(err) {
        console.error("Error fetching session data: ", err);
    }

    cardsSpace.forEach(card => {
        card.addEventListener("click", ()=> {
            const spaceID = card.getAttribute("data-space-id");
            console.log("selected space id: ", spaceID);
                 
        })
    })

})