// remove flash message after 6 seconds
const flash = document.querySelector('body div.flash')
if (flash) {
    setTimeout(() => {
        flash.style.display = 'none'
    }, 3000);
}

