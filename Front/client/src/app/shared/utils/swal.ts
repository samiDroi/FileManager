import SwalOriginal from 'sweetalert2';

const Swal = SwalOriginal.mixin({
  didOpen: () => {
    const container = SwalOriginal.getContainer();
    if (container) {
      container.style.zIndex = '100000';
    }
  }
});

export default Swal;
