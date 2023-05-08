import React, {useRef, useState} from 'react'
import Consola from '@monaco-editor/react'
import Salida from './Salida'
import Service from '../../Services/Service'

const Editor = () => {
    const [myValue, setMyValue] = useState('');
    //Lee archivo entrada
    const [ response, setResponse ] = useState("")
    const editorRef = useRef(null);
    const readFile = (e) => {
        const file = e.target.files[0];
        if(!file) return;

        const fileReader = new FileReader();
        fileReader.readAsText(file);

        fileReader.onload = () => {
            console.log(fileReader.result);
            setMyValue(fileReader.result);
        }

        fileReader.onerror = () => {
            console.log(fileReader.error);
        }
    }

    //Limpiar consola
    const handlerClear = () => {
        setResponse("")
    }

    //Corre el script
    const handleSave = () => {
    
        Service.parse(editorRef.current.getValue())
        .then((consola) => {
            setResponse(consola.result);
        
        })

    }

    const handleEditorDidMount = ( editor, monaco ) => {
        editorRef.current = editor;
    }

    const [contentMarkdown, setContentMarkdown] = useState('');

    return(
        <div>
            <div>
                <input type="file" multiple={false} onChange={readFile}/>
                <div className='btn-group'>
                    <button type="button" className="btn btn-success" onClick={handleSave}>Ejecutar</button>
                    <button type="button" className="btn btn-warning" onClick={handlerClear}>Limpiar</button>
                </div>
            </div>
            <Consola
                height = '450px'
                width = '99vw'
                theme = 'vs-dark'
                defaultLanguage='shell'
                value= { myValue }
                onChange={(value) => setContentMarkdown(value)}
                onMount={handleEditorDidMount}
            />
            <h3>Salida</h3>
            <Salida value={response} rows={10}/>
        </div>
    )
}

export default Editor