import React from "react";

function Salida(props) {

    const handlerChangeEditor = (evt) => {
        props.handlerChange(evt.target.value)
    }

    return (
        <>
            <div class="form-group editor-react">
                <label for="exampleFormControlTextarea1" style={{fontSize: "50px"}}></label>
                <textarea class="form-control" id="exampleFormControlTextarea1" rows={props.rows} onChange={handlerChangeEditor} value={props.value}></textarea>
                {props.comp}
            </div>
            <div id="lienzo">

            </div>
        </>
    )
}

export default Salida;