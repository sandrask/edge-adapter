<!--
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
-->

<template>
    <div class="hello">
        <h1>Please provide credentials</h1>
        <!--
            TODO pretty UI
        -->

        <textarea v-model="presentationRequest"></textarea>
    </div>
</template>

<script>
    export default {
        name: 'credentials',
        created: async function() {
            await this.$polyfill.loadOnce()
            this.getRequestForPresentation()
            const credentialQuery = {
                web: {
                    VerifiablePresentation: {
                        query: [
                            {
                                type: "presentationDefinitionQuery",
                                presentationDefinitionQuery: this.presentationRequest.pd

                            },
                            {
                                type: "DIDComm",
                                invitation: this.presentationRequest.invitation
                            }
                        ]
                    }
                }
            }
            const webCredential = await navigator.credentials.get(credentialQuery)
            console.log("received from user: " + JSON.stringify(webCredential))
            const redirectURL = this.validatePresentation(webCredential)
            // redirect user
            console.log(`redirecting user to ${redirectURL}`)
            // TODO add redirection
        },
        data() {
            return {
                presentationRequest: null
            }
        },
        methods: {
            getRequestForPresentation() {
                const handle = this.$route.query.pd
                console.info(`using handle: ${handle}`)
                this.$http.get(`/presentations/create?pd=${handle}`).then(
                    resp => {
                        this.presentationRequest = JSON.stringify(resp.data, null, 2)
                        console.log(`exchanged handle=${handle} for a request=${resp}`)
                    },
                    err => {
                        console.error(`failed to retrieve presentation-definitions: ${err}`)
                    }
                )
            },
            validatePresentation(presentation) {
                // TODO show error page if validation fails: https://github.com/trustbloc/edge-adapter/issues/128
                let redirectURL = ""
                this.$http.post(`/presentations/handleResponse`).then(
                    resp => {
                        redirectURL = resp.data.redirectURL
                        console.log(`received redirect url: ${redirectURL}`)
                    },
                    err => {
                        console.error(`failed to validate chapi response: ${err}`)
                    }
                )
                console.log(`submitted presentation=${presentation} and got redirectURL=${redirectURL}`)
                return redirectURL
            }
        }
    }
</script>
