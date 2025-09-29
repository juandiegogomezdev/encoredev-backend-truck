document.getElementById('registerForm').addEventListener('submit', async function (e) {
  const formContainer = document.getElementById('formContainer')
  const successfulContainer = document.getElementById('successfulContainer')
  const errorMessage = document.getElementsByClassName('errorMessage')[0]
  e.preventDefault()

  const params = new URLSearchParams(window.location.search)
  const token = params.get('token')

  const email = document.getElementById('email').value

  try {
    const response = await fetch(window.APP_CONFIG.url_register, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ email, token })

    })
    if (!response.ok) {
      errorMessage.textContent = await response.text()
      errorMessage.style.display = 'block'
      return
    }
    else {
      errorMessage.style.display = 'none'
      formContainer.style.display = 'none'
      successfulContainer.style.display = 'block'
    }

    
  } catch {
    errorMessage.style.display = 'block'
  }
})
