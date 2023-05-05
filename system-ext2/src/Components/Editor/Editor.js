import React from 'react'
import Consola from '@monaco-editor/react'

const Editor = () => {
    return(
        <div>
            <div>
                <input type="file" name="file"/>
                <div className='btn-group'>
                    <button type="button" className="btn btn-info">Cargar</button>
                    <button type="button" className="btn btn-success">Ejecutar</button>
                    <button type="button" className="btn btn-warning">Limpiar</button>
                </div>
            </div>
            <Consola
                height = '450px'
                width = '99vw'
                theme = 'vs-dark'
            />
            <h3>Salida</h3>
            <Consola
                height = '450px'
                width = '99vw'
                theme = 'vs-dark'
            />
        </div>
    )
}

export default Editor