function getData() {
    fetch('http://127.0.0.1:8080/api/crosshairs')
    .then(response => response.json())
			.then(data => {
				// Get the table element
				const table = document.getElementById("myTableRows");

				// Iterate over the data and add a row for each object
				data.forEach(item => {
					// Parse the DateTime string into a Date object
					const date = new Date(item.DateTime);
					// Format the date and time in a more readable format
					const formattedDateTime = date.toLocaleString();
					const steamURL = ['https://steamcommunity.com/profiles/' ,item.Steamid].join('')
					const row = document.createElement("tr");
					row.innerHTML = `
						<td>${formattedDateTime}</td>
						<td> 
							<a href="${steamURL}" rel="noopener" target="_blank">${item.Steamid}</a>
						</td>
						<td>${item.Playername}</td>
						<td style="font-family:Courier New;">${item.Crosshaircode}</td>
						<td>${item.Demoname}</td>
					`;
					table.appendChild(row);
				});
			});
  }

getData();