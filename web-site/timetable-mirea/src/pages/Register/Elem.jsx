import r from './Register.module.css';

const Elem = (props) => {
    return (
        <div className={r.form_reg}>
        <p className={r.text_form}>{props.name} </p><input className={r.form_input} type={props.type}/>
        </div>
    );
}

export default Elem;