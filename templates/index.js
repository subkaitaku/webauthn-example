// Base64 to ArrayBuffer
const bufferDecode = (value) =>
  Uint8Array.from(atob(value), (c) => c.charCodeAt(0));

// ArrayBuffer to URLBase64
const bufferEncode = (value) =>
  btoa(String.fromCharCode(...new Uint8Array(value)))
    .replace(/\+/g, "-")
    .replace(/\//g, "_")
    .replace(/=/g, "");

const registerUser = () => {
  const username = document.getElementById("email").value;
  if (username === "") {
    alert("Please enter a username");
    return;
  }

  fetch(`/register/begin/${username}`)
    .then((response) => response.json())
    .then((credentialCreationOptions) => {
      credentialCreationOptions.publicKey.challenge = bufferDecode(
        credentialCreationOptions.publicKey.challenge
      );
      credentialCreationOptions.publicKey.user.id = bufferDecode(
        credentialCreationOptions.publicKey.user.id
      );
      if (credentialCreationOptions.publicKey.excludeCredentials) {
        credentialCreationOptions.publicKey.excludeCredentials.forEach(
          (item) => {
            item.id = bufferDecode(item.id);
          }
        );
      }

      return navigator.credentials.create({
        publicKey: credentialCreationOptions.publicKey,
      });
    })
    .then((credential) => {
      const attestationObject = credential.response.attestationObject;
      const clientDataJSON = credential.response.clientDataJSON;
      const rawId = credential.rawId;

      return fetch(`/register/finish/${username}`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          id: credential.id,
          rawId: bufferEncode(rawId),
          type: credential.type,
          response: {
            attestationObject: bufferEncode(attestationObject),
            clientDataJSON: bufferEncode(clientDataJSON),
          },
        }),
      });
    })
    .then((response) => response.json())
    .then(() => {
      alert("Successfully registered " + username + "!");
    })
    .catch((error) => {
      console.log(error);
      alert("Failed to register " + username);
    });
};

const loginUser = () => {
  const username = document.getElementById("email").value;
  if (username === "") {
    alert("Please enter a username");
    return;
  }

  fetch("/login/begin/" + username)
    .then((response) => response.json())
    .then((credentialRequestOptions) => {
      credentialRequestOptions.publicKey.challenge = bufferDecode(
        credentialRequestOptions.publicKey.challenge
      );
      credentialRequestOptions.publicKey.allowCredentials.forEach((item) => {
        item.id = bufferDecode(item.id);
      });

      return navigator.credentials.get({
        publicKey: credentialRequestOptions.publicKey,
      });
    })
    .then((assertion) => {
      const authData = assertion.response.authenticatorData;
      const clientDataJSON = assertion.response.clientDataJSON;
      const rawId = assertion.rawId;
      const sig = assertion.response.signature;
      const userHandle = assertion.response.userHandle;

      return fetch("/login/finish/" + username, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          id: assertion.id,
          rawId: bufferEncode(rawId),
          type: assertion.type,
          response: {
            authenticatorData: bufferEncode(authData),
            clientDataJSON: bufferEncode(clientDataJSON),
            signature: bufferEncode(sig),
            userHandle: bufferEncode(userHandle),
          },
        }),
      });
    })
    .then((response) => response.json())
    .then(() => {
      alert("Successfully logged in " + username + "!");
    })
    .catch((error) => {
      console.log(error);
      alert("Failed to login " + username);
    });
};
