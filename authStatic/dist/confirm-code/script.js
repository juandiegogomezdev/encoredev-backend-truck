const params = new URLSearchParams(window.location.search)
const token = params.get('token')

document.addEventListener('DOMContentLoaded', () => {
  const confirmForm = document.getElementById('confirmForm')
  const inputCode = document.getElementById('inputCode')
  const errorMessage = document.getElementById('errorMessage')

  confirmForm.addEventListener('submit', async function (e) {
    e.preventDefault()

    console.log("Se envia el formulario")

    const errorMessage = document.getElementById('errorMessage')
    const code = inputCode.value.trim()

    // Check code length
    if (code.length != 6) {
      errorMessage.textContent = 'El codigo debe tener 6 numeros.'
      errorMessage.style.display = 'block'
      return
    }

    // Check only numbers
    if (!/^\d{6}$/.test(code)) {
      errorMessage.textContent = 'El codigo debe contener solo numeros.'
      errorMessage.style.display = 'block'
      return
    }
    console.log("Codigo validado, se envia al servidor")
    try {

      const response = await fetch(window.APP_CONFIG.url_confirm_code, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ code })
      })

      if (!response.ok) {
        errorMessage.textContent = await response.text()
        errorMessage.style.display = 'block'
        return
      } else {
        window.location.href = window.APP_CONFIG.page_org_select
      }
    } catch (error) {
      errorMessage.textContent = 'Error en el servidor!'
      errorMessage.style.display = 'block'
      return
    }
  })
})