<?php
require_once "../conexion.php";
session_start();
$id_user = $_SESSION['idUser'];
$eliminar = mysqli_query($conexion, "DELETE FROM graduaciones_temp WHERE id_usuario = $id_user");
if ($eliminar) {
    //echo "<div class='alert alert-danger'>Graduacion Borrada Correctamente</div>";
    echo
    "<script>Swal.fire({
        position: 'top-mid',
        icon: 'success',
        title: 'Graduacion Borrada Correctamente',
        showConfirmButton: false,
        timer: 2000
    })
    </script>;";   
    echo'
            <script type="text/javascript">;
            
            document.getElementById("borrar_grad").setAttribute("type", "hidden");
            
            </>';            
}else {
    //echo "<script>alert('Error al borrar Graduacion')</>";
    echo
    "<script>Swal.fire({
        position: 'top-mid',
        icon: 'error',
        title: 'Error al borrar Graduacion',
        showConfirmButton: false,
        timer: 2000
    })
    </script>;"; 
}

?>

