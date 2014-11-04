(function() {
	window.u2f_enroll = u2f_enroll;
	window.u2f_sign = u2f_sign;

	function ajax(url, args, cb) {
		var aj = new XMLHttpRequest();

		aj.onreadystatechange = function() {
			if(aj.readyState == 4 && aj.status == 200) {
				cb(JSON.parse(aj.responseText));
			} else {
				msg("failed: "+ aj.responseText);
			}
		}

		aj.open('POST', url, true);
		aj.setRequestHeader('Content-type', 'text/json')
		aj.send(JSON.stringify(args));
	}

	function u2f_enroll() {
		ajax("/enroll", {}, function(r) {
			msg("touch it to enroll");

			u2f.register(
				[r],
				[],
				function(response) {
					if(response.errorCode) {
						msg("failed to enroll:" + response.errorCode);
						return
					}
					msg("binding...")
					u2f_bind(response);
				}
			);
		});
	}

	function u2f_bind(enroll) {
		ajax("/bind", enroll, function(r) {
			msg("enrolled");
		});
	}

	function u2f_sign() {
		ajax("/sign", {}, function(r) {
			msg("touch it to sign in");
			console.log('sign', r);

			u2f.sign(r, function(response) {
				console.log(response);
				if (response.errorCode) {
					msg(response.errorCode);
					return;
				}

				msg("verifying");
				u2f_verify(response);
			}, 5);
		});
	}

	function u2f_verify(verify) {
		ajax("/verify", verify, function() {
			msg("logged in");
		});
	}

	function msg(m) {
		var e = document.getElementById("msg");
		e.innerHTML = 'HI: ';
		e.innerText = m;
	}

}());
