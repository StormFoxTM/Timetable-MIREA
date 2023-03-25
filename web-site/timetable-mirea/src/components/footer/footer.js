import './footer.css'

function Footer(props) {
    return(
    <footer>
        <p className='footer_main'>{props.info}</p>
    </footer>
    );
}

export default Footer;