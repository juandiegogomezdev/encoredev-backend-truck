

addEventListener("DOMContentLoaded", async (event) => {
    // List of memberships
    var membershipsList = [];

    const cardsSpace = document.querySelectorAll(".cardSpace");

    membershipsList = await loadMembershipCards()

    // cardsSpace.forEach(card => {
    //     card.addEventListener("click", ()=> {
    //         const spaceID = card.getAttribute("data-space-id");
    //         console.log("selected space id: ", spaceID);
                 
    //     })
    // })

})



async function loadMembershipCards() {
    
    // Memberships list
    var membershipsList = [];


    try {
        // Fetch memberships
        const response = await fetch(`${window.APP_CONFIG.api_url}/memberships`)
        if (response.ok) {
            const res = await response.json();
            membershipsList = res.memberships;

            console.log("Memberships fetched: ", membershipsList);

            // If no memberships, create personal membership
            if (membershipsList.length == 0) {
                membershipsList = await CreatePersonalMembership();
            }
        }
        else {
            console.error("Error fetching memberships: ", response.status);
        }  
    }
    catch (err) {
        console.error("Error fetching memberships: ", err);
    }

    return membershipsList;
}

async function CreatePersonalMembership() {
    console.log("Creating personal membership")
    var membershipsList = [];
    try {
        const response = await fetch(`${window.APP_CONFIG.api_url}/org/personal`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            }
        })
        if (response.ok) {
            const res = await response.json();
            membershipsList = [res.membership];
        }

    }
    catch (err) {
        console.error("Error creating personal membership card:", err)
    }

    return membershipsList;
}