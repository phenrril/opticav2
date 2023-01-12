<?php 
require "../conexion.php";
session_start();


$id_venta = $_POST['idventa'];
$id_abona = $_POST['idabona'];
if($id_venta == "" || $id_abona == ""){
    echo "<br><br><br><div class='row justify-content-center'><div class='alert alert-danger w-20'><div class='col-md-12 text-center'>COMPLETE AMBOS CAMPOS</div></div></div>";
    exit;
}
$query = mysqli_query($conexion, "SELECT * FROM postpagos WHERE id_venta = $id_venta");



if (mysqli_num_rows($query) > 0) {
$valueventa = mysqli_fetch_assoc($query);
$id_cliente = $valueventa['id_cliente'];
$abonatabla = $valueventa['abona'];
$abonatotal = $abonatabla + $id_abona;
$resto = $valueventa['resto'];
$precio = $valueventa['precio'];
$resto = $resto - $id_abona;
$update = mysqli_query($conexion, "UPDATE postpagos SET abona = '".$abonatotal."', resto = '".$resto."' WHERE id_venta = '".$id_venta."'");
if($update){
    $result = mysqli_affected_rows($conexion);
    if($result > 0){
        echo "<br><br><br><div class='row justify-content-center'><div class='alert alert-success w-20'><div class='col-md-12 text-center'>ID AGREGADO, VER PDF</div></div></div>";
        echo "  <div class='row justify-content-center'>
                    <a href='pdf/generar.php?cl=$id_cliente&v=$id_venta' target='_blank' class='btn btn-danger'><i class='fas fa-file-pdf'></i></a>
                <div>";    
    }
        else{
            echo "<br><br><br><div class='row justify-content-center'><div class='alert alert-danger w-30'><div class='col-md-12 text-center'>ERROR ACTUALIZANDO VENTA</div></div></div>";
        }
    }
}
else{
    echo "<br><br><br><div class='row justify-content-center'><div class='alert alert-danger w-30'><div class='col-md-12 text-center'>VENTA INEXISTENTE</div></div></div>";
}

?>
