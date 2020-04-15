var appCreateAccount = Vue.component("app-create-account", {
    template: `
	<div class="row align-items-center">
        <div class="col-sm-1 col-md-2 col-lg-3"></div>
		<div class="col-sm-10 col-md-8 col-lg-6 v-middle">
			<div class="login-panel">
				<form>
					<div class="form-group row">
						<label for="login-username" class="col-sm-12 col-form-label">Username:</label>
						<div class="col-sm-12">
							<input
								type="text"
								class="form-control"
								id="login-username"
                                spellcheck="false"
                                autocomplete="off"
                                v-model="Username"
                                name="username"
                                v-validate="'required'"
                                v-bind:class="{'form-control': true, 'error': errors.has('username') }"
								required
                            >
                            <div v-show="errors.has('username')" class="form-error">{{ errors.first('username') }}</div>
						</div>
					</div>
					<div class="form-group row">
						<label for="login-password" class="col-sm-12 col-form-label">Password:</label>
						<div class="col-sm-12">
							<input
								type="password"
								class="form-control"
								id="login-password"
								spellcheck="false"
								v-model="Password"
                                name="password"
                                v-validate="'required|min:6'"
                                ref="password"
                                v-bind:class="{'form-control': true, 'error': errors.has('password') }"
								required
							>
                            <div v-show="errors.has('password')" class="form-error">{{ errors.first('password') }}</div>
						</div>
					</div>
					<div class="form-group row">
						<label for="login-retype-password" class="col-sm-12 col-form-label">Retype Password:</label>
						<div class="col-sm-12">
							<input
								type="password"
								class="form-control"
								id="login-retype-password"
								spellcheck="false"
								v-model="RetypePassword"
                                name="password-retype"
                                v-validate="'required|confirmed:password'"
                                v-bind:class="{'form-control': true, 'error': errors.has('password-retype') }"
								required
							>
                            <div v-show="errors.has('password-retype')" class="form-error">{{ errors.first('password-retype') }}</div>
						</div>
					</div>
                    <hr>
					<div class="form-group row">
						<div class="col-sm-12">
							<button
								class="btn btn-primary btn-lg btn-login"
                                type="submit"
                                :disabled="errors.any() || !isComplete"
                                @click.prevent="createAccount()"
							>Create Account</button>
						</div>
					</div>
				</form>
			</div>
		</div>
    </div>
    `,
    $_veeValidate: {
        validator: "new"
    },
	data() {
		return {
            url: Config.Hostname + Config.AdminDir + "/" + Config.ApiPath,
			Username: "",
			Password: "",
			RetypePassword: ""
		};
    },
    computed: {
        isComplete () {
            return this.Username && this.Password && this.RetypePassword;
        }
    },
    methods: {
        createAccount() {
            if (this.Username == "" || this.Password == "" || this.RetypePassword == "") {
                return
            }
            if (this.Password != this.RetypePassword) {
                // highlight password don't match
                return
            }

			axios
				.post(
					this.url + "/create_account",
					{
                        username: this.Username,
                        password: this.Password,
					},
					{
						headers: {
							"content-type": "application/json"
						}
					}
				)
				.then(response => {
                    console.log(response);
                    this.mainBus.$emit("createdAccount");
				})
				.catch(error => {
					console.log(error);
				});
        }
    }
})