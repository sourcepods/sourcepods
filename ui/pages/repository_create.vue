<template>
	<div class="uk-container uk-container-small">
		<form class="uk-form-stacked" v-on:submit.prevent="create">

			<div class="uk-margin">
				<label class="uk-form-label" for="name">Name</label>
				<div class="uk-form-controls">
					<input class="uk-input" id="name" type="text" autofocus required v-model="repository.name">
				</div>
			</div>

			<div class="uk-margin">
				<label class="uk-form-label" for="description">Description</label>
				<div class="uk-form-controls">
					<input class="uk-input" id="description" type="text" v-model="repository.description">
				</div>
			</div>

			<div class="uk-margin">
				<label class="uk-form-label" for="website">Website</label>
				<div class="uk-form-controls">
					<input class="uk-input" id="website" type="text" v-model="repository.website">
				</div>
			</div>

			<div class="uk-margin uk-grid-small uk-child-width-auto" uk-grid>
				<label><input class="uk-radio" type="radio" name="radio2" :value="false" v-model="repository.private"
							  checked>
					Public</label>
				<label><input class="uk-radio" type="radio" name="radio2" :value="true" v-model="repository.private">
					Private</label>
			</div>

			<div class="uk-margin">
				<button class="uk-button uk-button-primary" :disabled="loading" type="submit">Create Repository</button>
			</div>
		</form>

	</div>
</template>

<script>
	export default {
		data() {
			return {
				loading: false,
				repository: {}
			}
		},
		methods: {
			setLoading(isLoading) {
				this.loading = isLoading;
				this.$store.commit('loading', isLoading)
			},
			create(e) {
				this.setLoading(true);

				this.$store.dispatch('createRepository', this.repository)
					.then((res) => {
						this.$router.push(`/${res.owner.username}/${res.name}`);
						this.setLoading(false);
						UIkit.notification('Repository created', {status: 'success', pos: 'bottom-center'});
					})
					.catch((err) => {
						this.setLoading(false);
						UIkit.notification('Repository not created', {status: 'danger', pos: 'bottom-center'})
					})
			}
		}
	}
</script>
