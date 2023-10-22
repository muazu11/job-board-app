import { env } from '$env/dynamic/public'

const baseRoute = "http://" + env.PUBLIC_API_HOST + ":" + env.PUBLIC_API_PORT + "/"

export function Advertisement(id, title, description, wage, address, zipCode, city, workTime, companyName, companySiren, companyLogoURL, applied) {
  this.id = id
  this.title = title
  this.description = description
  this.wage = wage
  this.address = address
  this.zipCode = zipCode
  this.city = city
  this.workTime = workTime
  this.companyName = companyName
  this.companySiren = companySiren
  this.companyLogoURL = companyLogoURL
  this.applied = applied
}

export async function getAllAds(token = "", pageCursor = 0, pagePrevious = false) {
  let url = baseRoute + 'advertisements/with_detail?'+ new URLSearchParams({
    pageCursor: pageCursor,
    pagePrevious: pagePrevious
  })
  let toFetch = fetch(url)
  if (token !== "" && token !== "undefined") {
    toFetch = fetch(url, {
      headers: {
        "Authorization": "Basic " + token
      },
      body: JSON.stringify({
        pageCursor: pageCursor,
        pagePrevious: pagePrevious
      })
    })
  }
  let promise = toFetch.then(response => (response.json()))
    .then((rep) => {
      let advertisements = []
      rep["Data"].forEach((jsonAdvertisement) => {
        advertisements.push(new Advertisement(
          jsonAdvertisement.ID, jsonAdvertisement.Title,
          jsonAdvertisement.Description, jsonAdvertisement.Wage,
          jsonAdvertisement.Address, jsonAdvertisement.ZipCode, jsonAdvertisement.City,
          jsonAdvertisement.WorkTimeNs / 3600000000000,
          jsonAdvertisement["Company"].Name, jsonAdvertisement["Company"].Siren,
          jsonAdvertisement["Company"].LogoURL, jsonAdvertisement.Applied))
      })
      return [advertisements, rep.Cursors.Previous, rep.Cursors.Next]
    })

  return await promise


}

export async function submitApply(message, advertisement_ID, token) {
  let url = baseRoute + 'applications/me?' + new URLSearchParams({
    message: message,
    advertisement_id: advertisement_ID,
  })

  let promise = await fetch(url, {
    method: 'POST',
    headers: {
      "Authorization": "Basic " + token
    },
    mode: "cors",
  })
}

export async function getSvg(fileName) {
  let promise = fetch(fileName)
    .then((res) => res.text())
    .then((text) => {
      return text;
    })
    .catch((e) => console.error(e));
  return await promise;
}

export async function getMe(token) {
  let url = baseRoute + 'users/me'
  let promise = fetch(url, {
    method: 'GET',
    headers: {
      "Authorization": "Basic " + token
    }
  })
    .then(response => (response.json()))
    .then((data) => {
      return data
    })
    .catch((error) => {
      return false
    })
  return await promise
}

export function getCookie(cname) {
  let name = cname + "=";
  let decodedCookie = decodeURIComponent(document.cookie);
  let ca = decodedCookie.split(';');
  for (let i = 0; i < ca.length; i++) {
    let c = ca[i];
    while (c.charAt(0) === ' ') {
      c = c.substring(1);
    }
    if (c.indexOf(name) === 0) {
      return c.substring(name.length, c.length);
    }
  }
  return "";
}

export async function sendCredentials(login, password) {

  let url = baseRoute + 'users/login'
  let promise =
    fetch(url, {
      method: 'POST',
      body: JSON.stringify({
        email: login,
        password: password
      })
    })
      .then(response => (response.json())
        .then(data => {
          return data.Token
        }))
      .catch(error => {
        console.log(error)
      })
  return await promise
}

export async function createUser(email, password, name, surname, tel, birthDate) {
  let url = baseRoute + 'users'
  let promise = fetch(url, { method: 'POST'
  , body: JSON.stringify({
      email: email,
      name: name,
      surname: surname,
      phone: tel,
      dateOfBirthUTC: birthDate,
      password: password,
      role: "user"})
  })
    .then(response => {return response.status===201})
    .catch(error => {
      console.log(error)
    })
  return await promise
}

export async function createAdvertisement(title, description, companyID, address, city, zipCode,workTimeNs) {
  let url = baseRoute + 'advertisements'
  let promise = fetch(url, { method: 'POST'
    , body: JSON.stringify({
        title: title,
        description: description,
        company_id: companyID,
        address: address,
        city: city,
        zip_code: zipCode,
        work_time_ns: workTimeNs,
    })
  })
      .then(response => {return response.status===201})
      .catch(error => {
        console.log(error)
      })
  return await promise
}

export async function createApplication(advertisementID, applicantID, message) {
  let url = baseRoute + 'advertisements'
  let promise = fetch(url, { method: 'POST'
    , body: JSON.stringify({
        advertisement_id: advertisementID,
        applicant_id: applicantID,
        message: message,
    })
  })
      .then(response => {return response.status===201})
      .catch(error => {
        console.log(error)
      })
  return await promise
}

export async function createCompany(name, logoURL, siren) {
  let url = baseRoute + 'advertisements'
  let promise = fetch(url, { method: 'POST'
    , body: JSON.stringify({
        name: name,
        logoURL: logoURL,
        siren: siren,
    })
  })
      .then(response => {return response.status===201})
      .catch(error => {
        console.log(error)
      })
  return await promise
}
export async function getAllApplications(token, pageCursor = 0, pagePrevious = false){
  let url = baseRoute + 'applications?'+ new URLSearchParams({
    pageCursor: pageCursor,
    pagePrevious: pagePrevious
  })
  let promise = fetch(url, { method: 'GET',headers: {
      "Authorization": "Basic " + token
    }
  })
      .then(response => (response => (response.json())))
      .then((data) => {
        return data
      })
      .catch(error => {
        console.log(error)
      })
  return await promise
}
export async function getAllUsers(token, pageCursor = 0, pagePrevious = false){
  let url = baseRoute + 'users?'+new URLSearchParams({
    pageCursor: pageCursor,
    pagePrevious: pagePrevious
  })
  let promise = fetch(url, { method: 'GET',headers: {
      "Authorization": "Basic " + token
    }
  })
      .then(response => (response => (response.json())))
      .then((data) => {
        return data
      })
      .catch(error => {
        console.log(error)
      })
  return await promise
}


export async function getAllcompanies(token, pageCursor = 0, pagePrevious = false){
  let url = baseRoute + 'companies?'+new URLSearchParams({
    pageCursor: pageCursor,
    pagePrevious: pagePrevious
  })
  let promise = fetch(url, { method: 'GET',headers: {
      "Authorization": "Basic " + token
    }
  })
      .then(response => (response => (response.json())))
      .then((data) => {
            return data
          })
      .catch(error => {
        console.log(error)
      })
  return await promise
}

export async function updateAdvertisement(id, title, description, wage, address, zipCode, city, workTime, companyName, companySiren, companyLogoURL, applied, token) {
    let url = baseRoute + 'advertisements/' + id
    let promise = fetch(url, {
        method: 'PUT',
        headers: {
        "Authorization": "Basic " + token
        },
        body: JSON.stringify({
          title: title,
          description: description,
          wage: wage,
          address: address,
          zip_code: zipCode,
          city: city,
          work_time_ns: workTime,
          company_name: companyName,
          company_siren: companySiren,
          company_logo_url: companyLogoURL,
          applied: applied
        })
    })
        .then(data => {
            return data.status
        })
        .catch(error => {
            return false
        })
    return await promise
}

export async function updateApplication(id, message, advertisementID, applicantID, token) {
  let url = baseRoute + 'applications/' + id
  let promise = fetch(url, {
    method: 'PUT',
    headers: {
      "Authorization": "Basic " + token
    },
    body: JSON.stringify({
        message: message,
        advertisement_id: advertisementID,
        applicant_id: applicantID,
    })
  })
      .then(data => {
        return data.status
      })
      .catch(error => {
        return false
      })
  return await promise
}

export async function updateCompany(id, name, logoURL, siren, token) {
  let url = baseRoute + 'advertisements/' + id
  let promise = fetch(url, {
    method: 'PUT',
    headers: {
      "Authorization": "Basic " + token
    },
    body: JSON.stringify({
        name: name,
        logoURL: logoURL,
        siren: siren,
    })
  })
      .then(data => {
        return data.status
      })
      .catch(error => {
        return false
      })
  return await promise
}

export async function updateUser(id, email, name, surname, tel, birthDate, token) {
  let url = baseRoute + 'advertisements/' + id
  let promise = fetch(url, {
    method: 'PUT',
    headers: {
      "Authorization": "Basic " + token
    },
    body: JSON.stringify({
        email: email,
        name: name,
        surname: surname,
        phone: tel,
        dateOfBirthUTC: birthDate,
        role: "user"
    })
  })
      .then(data => {
        return true
      })
      .catch(error => {
        return false
      })
  return await promise
}

export async function deleteAdvertisement(id, token) {
  let url = baseRoute + 'advertisements/' + id
    let promise = fetch(url, {
        method: 'DELETE',
        headers: {
          "Authorization": "Basic " + token
        }
    })
}
export function deleteApplication(id, token) {
  let url = baseRoute + 'applications/' + id
  let promise = fetch(url, {
    method: 'DELETE',
    headers: {
      "Authorization": "Basic " + token
    }
  })
}

export function deleteCompany(id, token) {
  let url = baseRoute + 'companies/' + id
  let promise = fetch(url, {
    method: 'DELETE',
    headers: {
      "Authorization": "Basic " + token
    }
  })
}

export function deleteUser(id, token) {
  let url = baseRoute + 'users/' + id
  let promise = fetch(url, {
    method: 'DELETE',
    headers: {
      "Authorization": "Basic " + token
    }
  })
}

async function updatePwd(password, token) {
  let url = baseRoute + 'users/password/me?' + new URLSearchParams({
    password: password,
  })
  let promise = fetch(url, {
    method: 'PUT',
    headers: {
      "Authorization": "Basic " + token
    }
  })
    .then(data => {
      return true
    })
    .catch(error => {
      return false
    })
  return await promise
}

export async function updateMyInfo(email, name, surname, tel, birthDate, token) {
  let url = baseRoute + 'users/me?' + new URLSearchParams({
    email: email,
    name: name,
    surname: surname,
    phone: tel,
    date_of_birth_utc: birthDate,
    role: "user"
  })
  let promise = fetch(url, {
    method: 'PUT',
    headers: {
      "Authorization": "Basic " + token
    }
  })
    .then(data => {
      return true
    })
    .catch(error => {
      return false
    })
  return await promise
}

export async function updateProfile(email, name, surname, tel, birthDate, token, password) {
  let infoOk = await updateMyInfo(email, name, surname, tel, birthDate, token)
  if (infoOk) {
    return await updatePwd(password, token)
  }
  return false
}
