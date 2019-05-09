function auth() {
    var currentUser = {};
    firebase.auth().onAuthStateChanged(function(user) {
        if (user) {
            // User is signed in.
            var displayName = user.displayName;
            var email = user.email;
            var emailVerified = user.emailVerified;
            var photoURL = user.photoURL;
            var isAnonymous = user.isAnonymous;
            var uid = user.uid;
            var providerData = user.providerData;
            currentUser.displayName = displayName;
            currentUser.uid = uid;
            // document.getElementById("login-name").textContent = displayName;
            // document.getElementById("login-photo-url").src = photoURL;

            firebase.auth().currentUser.getIdToken(/* forceRefresh */ true).then(function(idToken) {
                fetch("/api/hoge", {
                    method: "GET",
                    headers: {
                        "Content-Type": "application/json; charset=utf-8",
                        "Authorization": "Bearer " + idToken
                    },
                }).then(response => response.text())
                    .then(text => {
                        console.log(text);
                    });
            }).catch(function(error) {
                // Handle error
            });
        } else {
            // FirebaseUI config.
            var uiConfig = {
                signInSuccessUrl: '/',
                signInOptions: [
                    // Leave the lines as is for the providers you want to offer your users.
                    firebase.auth.GoogleAuthProvider.PROVIDER_ID,
                ],
                // tosUrl and privacyPolicyUrl accept either url string or a callback
                // function.
                // Terms of service url/callback.
                tosUrl: '<your-tos-url>',
                // Privacy policy url/callback.
                privacyPolicyUrl: function() {
                    window.location.assign('<your-privacy-policy-url>');
                }
            };
            // Initialize the FirebaseUI Widget using Firebase.
            var ui = new firebaseui.auth.AuthUI(firebase.auth());
            // The start method will wait until the DOM is loaded.
            ui.start('#firebaseui-auth-container', uiConfig);
        }
    });
}