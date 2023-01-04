function getData() {
    fetch('http://127.0.0.1:8080/api/crosshairs')
    .then(response => response.json())
			.then(data => {
				// Get the table element
				const table = document.getElementById("myTableRows");

				// Iterate over the data and add a row for each object
				data.forEach(item => {
					// Parse the DateTime string into a Date object
					const steamURL = ['https://steamcommunity.com/profiles/' ,item.steamid].join('')
					const row = document.createElement("tr");
					row.innerHTML = `
						<td> 
							<a href="${steamURL}" rel="noopener" target="_blank">${item.steamid}</a>
						</td>
						<td>${item.playername}</td>
						<td style="font-family:Courier New;">${item.matches[0].crosshair}</td>
						<td>${item.matches[0].name}</td>
					`;
					table.appendChild(row);
				});
			});
  }

getData();