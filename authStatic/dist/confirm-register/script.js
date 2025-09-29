const params = new URLSearchParams(window.location.search)
const token = params.get('token')

function validarContrasenia (pass) {
    return [
      pass.length >= 8, // Minimo 8 letras
      /[A-Z]/.test(pass), // Al menos una letrea mayuscula
      /[a-z]/.test(pass), // Al menos una letra minuscula
      /[\d]/.test(pass), // Al menos un numero
      /[\W_]/.test(pass)
    ]
  }


document.addEventListener('DOMContentLoaded', () => {

  const confirmForm = document.getElementById('confirmForm')
  const tryAgain = document.getElementById('tryAgain')
  const password1 = document.getElementById('password1')
  const password2 = document.getElementById('password2')

  
  const formContainer = document.getElementById('formContainer')
  const errorContainer = document.getElementById('errorContainer')
  const successContainer = document.getElementById('successContainer')


  const errorMessage1 = document.getElementsByClassName('errorMessage')[0]
  const errorMessage2 = document.getElementById('errorMessageFetch')

  confirmForm.addEventListener('submit', async function (e) {
    e.preventDefault()

    const pass1 = password1.value.trim()
    const pass2 = password2.value.trim()
    const validations = validarContrasenia(pass1)

    updateIcon('check-length', validations[0])
    updateIcon('check-uppercase', validations[1])
    updateIcon('check-lowercase', validations[2])
    updateIcon('check-number', validations[3])
    updateIcon('check-special', validations[4])

    const requirementsOk = validations.every(Boolean)

    if (!requirementsOk) {
      errorMessage1.textContent = 'La contrasenia no cumple todos los requisitos.'
      errorMessage1.style.display = 'block'
      return
    }

    if (pass1 !== pass2) {
      errorMessage1.textContent = 'Las contrasenias no coinciden.'
      errorMessage1.style.display = 'block'
      return
    }

    formContainer.style.display = 'none'

    try {
      const response = await fetch(window.APP_CONFIG.url_confirm, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ password: pass1, token })
      })

      if (!response.ok) {
        errorContainer.style.display = 'block'
        errorMessage2.textContent = await response.text()
        return
      }
      else {
        formContainer.style.display = 'none'
        successContainer.style.display = 'block'
        // redirect to org-select after 1 second
        setTimeout(() => {
          window.location.href = window.APP_CONFIG.page_org_select
        }, 1000)
      }


    } catch (error) {
      errorMessage2.style.display = 'block'
      errorMessage2.textContent = 'Oops, Ocurrio algun error!'
    }
  })

  tryAgain.addEventListener('click', () => {
    formContainer.style.display = 'block'
    errorContainer.style.display = 'none'
  })
})

function updateIcon (id, condition) {
  const icon = document.getElementById(id)
  icon.classList.remove('fa-circle-check', 'fa-circle-xmark')

  if (condition) {
    icon.classList.add('fa-circle-check')
    icon.style.color = 'lightgreen'
  } else {
    icon.classList.add('fa-circle-xmark')
    icon.style.color = 'red'
  }
}


