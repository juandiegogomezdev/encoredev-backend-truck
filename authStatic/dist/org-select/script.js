let membershipsList = null;
let filteredCompaniesList = null;
let currentFilter = null;
let loadingMemberships = true;
let countRefresingToken = 0;

addEventListener("DOMContentLoaded", async (event) => {

  spaceEventListener();
  const cardsSpace = document.querySelectorAll(".cardSpace");
  loadMembershipCards();
});


// Apply filter to membershipsList and update the UI
function applyFilter() {
  if (currentFilter != null && loadingMemberships) {
    showLoadingCard();
  } else {
    hiddenLoadingCard();
  }

  if (!membershipsList) return;


  // Filter companies based on spaceID
  if (currentFilter == "my-companies") {
    // filter membershipsList where org_type is 'company' and role is "owner" or "partner"
    filteredCompaniesList = membershipsList.filter(
      (membership) =>
        membership.org_type == "company" &&
        (membership.role_name == "owner" ||
          membership.role_name == "partner")
    );
  }
  
  // Filter jobs base on spaceID
  if (currentFilter == "my-jobs") {
    // filter membershipsList where org_type is 'company' and role is not "owner" or "partner"
    filteredCompaniesList = membershipsList.filter(
      (membership) =>
      membership.org_type == "company" &&
      membership.role_name != "owner" &&
      membership.role_name != "partner"
    );
  }
  if (currentFilter == "independent") {
    // Go to independent page
    filteredCompaniesList = membershipsList.filter(
      membershipsList => membershipsList.org_type == "personal" && membershipsList.role_name == "owner")
  
    }

  
  CreateMembershipCards()
}


// Add event listeners to space cards
async function spaceEventListener(membershipsList) {
  const cardsSpace = document.querySelectorAll(".cardSpace");
  const spaceSelectedTitle = document.querySelector(".spaceSelected");
  const spaceDescription = document.querySelector(".spaceDescription");

  cardsSpace.forEach((card) => {
    card.addEventListener("click", () => {
      currentFilter = card.getAttribute("data-space-id");
      applyFilter();
      toggleSelectedCardSpace(card);


      if (currentFilter == 'my-companies') {
        spaceSelectedTitle.textContent = "Mis empresas";
        spaceDescription.textContent = "Elije y administra tu empresa.";
      }
      if (currentFilter == 'my-jobs') {
        spaceSelectedTitle.textContent = "Mis trabajos";
        spaceDescription.textContent = "Selecciona el trabajo que deseas gestionar.";
      }

      if (currentFilter == 'independent') {
        spaceSelectedTitle.textContent = "Independiente";
        spaceDescription.textContent = "Gestiona tu perfil independiente.";
      }
    });
  });
}

// Show loading card
function showLoadingCard(card) {
  const loadingCard = document.getElementById("loadingCard");
  loadingCard.style.display = "block";
}

// HIde loading card
function hiddenLoadingCard(card) {
  const loadingCard = document.getElementById("loadingCard");
  loadingCard.style.display = "none";
}


// This refresh the token
async function refreshToken() {
  countRefresingToken++;
  if (countRefresingToken > 2) {
    throw new Error("Too many token refresh attempts");
  }
  try {
    const response = await fetch(`${window.APP_CONFIG.api_url}/auth/refresh`, {
      method: "POST",
      credentials: true,
    });
    if (!response.ok) {
      window.location.href = window.APP_CONFIG.base_url + '/static/login';
    }
  } catch (err) {
    console.log("Error refreshing token: ", err);
    throw err;
  } 
}

// Load memberships
async function loadMembershipCards() {

  try {
    // Fetch memberships
    const response = await fetch(`${window.APP_CONFIG.api_url}/memberships`);
    if (response.ok) {
      const res = await response.json();
      membershipsList = res.memberships;

      // If no memberships, create personal membership
      if (membershipsList.length == 0 && countRefresingToken < 2) {
        await CreatePersonalMembership();
      }

      loadingMemberships = false;
      applyFilter(currentFilter);

    } else if(response.status === 401){
      const res = await response.json();
      if (res.details == null) {
        window.location.href = window.APP_CONFIG.base_url + '/static/login';
      }
      else if (res.details.token_status == "token_expired") {
        // Token expired, try to refresh
        if (countRefresingToken <= 2) {
          await refreshToken();
          loadMembershipCards();

        }
      }
      else {

      }

    }
    
    else {
      // console.error("Error fetching memberships: ", response);
    }
  } catch (err) {
    console.error("Error fetching memberships: ", err);
  }
}


// Create personal membership
async function CreatePersonalMembership() {
  try {
    const response = await fetch(`${window.APP_CONFIG.api_url}/org/personal`, {
      method: "POST",
      credentials: "include",
    });
    if (response.ok) {
      const res = await response.json();
      membershipsList = [res.membership];

      loadingMemberships = false;
      applyFilter(currentFilter);
    } else if (response.status === 401) {
      if (countRefresingToken <= 2) {
        await refreshToken();
        await CreatePersonalMembership();
      }

    }
  } catch (err) {
    console.error("Error creating personal membership card:", err);
  }

}

const rolesClass = {
  'owner': 'rolOwner',
  'partner': 'rolPartner',
  'admin': 'rolAdmin',
  'driver': 'rolDriver',
  'accountant': 'rolAccountant'
}

const rolesTranslate = {
  'owner': 'Jefe',
  'partner': 'Socio',
  'admin': 'Administrador',
  'driver': 'Conductor',
  'accountant': 'Contador'
}

const statusClass = {
  'active': 'statusOrgActive',
  'revoked': 'statusOrgRevoked',
  'suspended': 'statusOrgSuspended',
  'finalized': 'statusOrgSuspended'

}

const statusTranslate = {
  'active': 'Activo',
  'revoked': 'Revocado',
  'suspended': 'Suspendido',
  'finalized': 'Finalizado'
}

async function CreateMembershipCards() {

  // Get container
  const cardsContainer = document.getElementById("companiesContainer");

  // Delete all cards
  const cards = document.querySelectorAll(".cardCompany");
  cards.forEach((card) => {
    if (card.id != "loadingCard") {
      card.remove();
    }
  });

  filteredCompaniesList.forEach((membership) => {
    const card = document.createElement("div");
    card.classList.add("card", "cardCompany");

    const cardHeader = document.createElement("div");
    cardHeader.classList.add("headerTopCompany");

    const cardRol = document.createElement("div");
    const cardRolClass = rolesClass[membership.role_name] || '';
    cardRol.classList.add("rol", cardRolClass);
    cardRol.textContent = rolesTranslate[membership.role_name] || '';

    const cardStatus = document.createElement("div");

    cardStatus.classList.add("statusOrg", statusClass[membership.status] || '');
    cardStatus.textContent = statusTranslate[membership.status] || '';

    const cardBody = document.createElement("div")
    cardBody.classList.add("cardCompanyBody");

    const orgName = document.createElement("div");
    orgName.classList.add("nameOrg");
    orgName.textContent = membership.org_name;

    const date = document.createElement("div");
    date.classList.add("date")
    date.textContent = `🗓️ ${convertoTimeStampToLegibleDate(membership.created_at)}`;
 
    cardHeader.appendChild(cardRol);
    cardHeader.appendChild(cardStatus);
    cardBody.appendChild(orgName);
    cardBody.appendChild(date);
    card.appendChild(cardHeader);
    card.appendChild(cardBody);

    // Add to container
    cardsContainer.appendChild(card);
  })

}



function convertoTimeStampToLegibleDate(timestamp) {
  // Convert timestamp to dd/mm/yyyy format
  const date = new Date(timestamp);
  const options = { day: '2-digit', month: '2-digit', year: 'numeric' };
  return date.toLocaleDateString(undefined, options);
}


// Toggle selected class on space cards
function toggleSelectedCardSpace(card) {
  const cardsSpace = document.querySelectorAll(".cardSpace");
  cardsSpace.forEach((c) => {
    c.classList.remove("selected");
  });

  // Check if card is already selected
  if (card.classList.contains("selected")) {
    card.classList.remove("selected");
    return;
  }
  card.classList.add("selected");
}
