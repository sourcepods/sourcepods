import Vue from "vue";
import ApolloClient, {createNetworkInterface} from "apollo-client";
import VueApollo from "vue-apollo";

export const apolloClient = new ApolloClient({
	networkInterface: createNetworkInterface({
		uri: `${window.config.api}/query`,
		transportBatching: true,
		opts: {
			credentials: 'same-origin',
		},
	}),
	connectToDevTools: true,
});


export function createProvider() {
	return new VueApollo({
		defaultClient: apolloClient,
	});
}

Vue.use(VueApollo);
