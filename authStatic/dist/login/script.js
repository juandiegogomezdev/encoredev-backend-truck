document.getElementById('loginForm').addEventListener('submit', async function (e) {
  e.preventDefault()

  const email = document.getElementById('email').value
  const password = document.getElementById('password').value
  const errorMessage = document.querySelector('.errorMessage')
  try {
    const response = await fetch(window.APP_CONFIG.url_login, {
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ email, password })

    })


    if (!response.ok) {
      errorMessage.textContent = await response.text()
      errorMessage.style.display = 'block'
    } else {
      window.location.href = window.APP_CONFIG.page_url_login_confirm
    }
  } catch (error) {
    console.log(error)
    errorMessage.style.display = 'block'
    errorMessage.textContent = 'Error en el servidor'
  }
})
