import React from 'react'

const Reportes = () => {

    function verRep(){
        var obj = { 'cmd': 'Reporte' }
            fetch(`http://18.191.10.110:5000/reports`, {
                method: 'POST',
                body: JSON.stringify(obj),
            })
                .then(res => res.json())
                .catch(err => {
                    console.error('Error:', err)
                    alert("Ocurrio un error, ver la consola")
                })
                .then(response => {
                    var image = new Image();
                    image.src = response.result
                    document.getElementById('RepDiv').appendChild(image);
                })
    }
    return(
        <div>
            <h3>Reportes</h3>
            <button type="submit" onClick={verRep} class="btn btn-primary" >Ver Reporte</button>
            <div class="modal-body">
                <div id="RepDiv"></div>
            </div>
        </div>
    )
}

export default Reportes