//-------------------------------------------------------------------
// Enrollment HFC Library
//-------------------------------------------------------------------

module.exports = function (logger) {
	var FabricClient = require('fabric-client');
	var path = require('path');
	var common = require(path.join(__dirname, './common.js'))(logger);
	var enrollment = {};
	var User = require('fabric-client/lib/User.js');
	var CaService = require('fabric-ca-client/lib/FabricCAClientImpl.js');
	var Orderer = require('fabric-client/lib/Orderer.js');
	var Peer = require('fabric-client/lib/Peer.js');
	FabricClient.setConfigSetting('request-timeout', 90000);

	//-----------------------------------------------------------------
	// Enroll with Admin Certs - use this for install || instantiate || creating a channel
	//-----------------------------------------------------------------
	/*
		options = {
			peer_urls: ['array of peer grpc urls'],
			channel_id: 'channel name',
			uuid: 'unique name for this enrollment',
			orderer_url: 'grpc://url_here:port',
			privateKeyPEM: '<cert here>',
			signedCertPEM: '<cert here>',
			msp_id: 'string',
			orderer_tls_opts: {
				pem: 'complete tls certificate',					<required if using ssl>
				common_name: 'common name used in pem certificate' 	<required if using ssl>
			},
			peer_tls_opts: {
				pem: 'complete tls certificate',					<required if using ssl>
				common_name: 'common name used in pem certificate' 	<required if using ssl>
			},
			kvs_path: '/path/to/the/key/value/store'
		}
	*/

	enrollment.enrollWithAdminCert = function (options, cb) {
		var client = new FabricClient();
		var channel = client.newChannel(options.channel_id);

		var debug = {														// this is just for console printing, no PEM here
			peer_urls: options.peer_urls,
			channel_id: options.channel_id,
			uuid: options.uuid,
			orderer_url: options.orderer_url,
			msp_id: options.msp_id,
		};
		logger.info('[fcw] Going to enroll with admin cert! ', debug);

		// Make eCert kvs (Key Value Store)
		FabricClient.newDefaultKeyValueStore({
			path: options.kvs_path 													//get eCert in the kvs directory
		}).then(function (store) {
			client.setStateStore(store);
			return getSubmitterWithAdminCert(client, options);						//admin cert is different
		}).then(function (submitter) {

			channel.addOrderer(new Orderer(options.orderer_url, {
				pem: options.orderer_tls_opts.pem,
				'ssl-target-name-override': options.orderer_tls_opts.common_name	//can be null if cert matches hostname
			}));

			channel.addPeer(new Peer(options.peer_urls[0], {						//add the first peer
				pem: options.peer_tls_opts.pem,
				'ssl-target-name-override': options.peer_tls_opts.common_name		//can be null if cert matches hostname
			}));
			logger.debug('added peer', options.peer_urls[0]);

			// --- Success --- //
			logger.debug('[fcw] Successfully got enrollment ' + options.uuid);
			if (cb) cb(null, { client: client, channel: channel, submitter: submitter });
			return;

		}).catch(function (err) {

			// --- Failure --- //
			logger.error('[fcw] Failed to get enrollment ' + options.uuid, err.stack ? err.stack : err);
			var formatted = common.format_error_msg(err);

			if (cb) cb(formatted);
			return;
		});
	};

	// Get Submitter - ripped this function off from helper.js in fabric-client
	function getSubmitterWithAdminCert(client, options) {
		return Promise.resolve(client.createUser({
			username: options.msp_id,
			mspid: options.msp_id,
			cryptoContent: {
				privateKeyPEM: options.privateKeyPEM,
				signedCertPEM: options.signedCertPEM
			}
		}));
	}

	return enrollment;
};
