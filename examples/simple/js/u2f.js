(function() {
	window.u2fRegister = u2fRegister;
	window.u2fSign = u2fSign;

	function ajax(url, args, cb) {
		var aj = new XMLHttpRequest();

		aj.onreadystatechange = function() {
			if (aj.readyState == 4) {
				if (aj.status == 200) {
					cb(JSON.parse(aj.responseText));
				} else {
					msg("failed: " + aj.responseText);
				}
			}
		}

		aj.open('POST', url, true);
		aj.setRequestHeader('Content-type', 'text/json')
		aj.send(JSON.stringify(args));
	}

	function u2fRegister() {
		ajax("/Register", {}, function(r) {
			msg("touch it to register");

			u2f.register(
				[r], [],
				function(response) {
					if (response.errorCode) {
						msg("failed to enroll:" + response.errorCode);
						return
					}
					msg("finalizing/validating registration...")
					u2fRegisterFin(response);
				}
			);
		});
	}

	function u2fRegisterFin(rr) {
		ajax("/RegisterFin", rr, function(r) {
			msg("device registered");
		});
	}

	function u2fSign() {
		ajax("/Sign", {}, function(r) {
			msg("touch it to sign in");
			console.log('sign', r);

			u2f.sign(r, function(response) {
				if (response.errorCode) {
					msg(response.errorCode);
					return;
				}

				msg("verifying");
				u2fSignFin(response);
			}, 5);
		});
	}

	function u2fSignFin(verify) {
		ajax("/SignFin", verify, function() {
			msg("logged in");
		});
	}

	function msg(m) {
		var e = document.getElementById("msg");
		e.innerHTML = 'HI: ';
		e.innerText = m;
	}

}());