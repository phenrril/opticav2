<?php 
require "../conexion.php";
session_start();
$id_venta = $_POST['idventa'];
$id_cristal = $_POST['idcristal'];
if($id_venta == "" || $id_cristal == ""){
    echo "<br><br><br><div class='row justify-content-center'><div class='alert alert-danger w-20'><div class='col-md-12 text-center'>COMPLETE AMBOS CAMPOS</div></div></div>";
    exit;
}
$query = mysqli_query($conexion, "SELECT * FROM ventas WHERE id = $id_venta");
$valueventa = mysqli_fetch_assoc($query);
//$id = $valueventa['id'];
$id_cliente = $valueventa['id_cliente'];

$update = mysqli_query($conexion, "UPDATE detalle_venta SET idcristal = $id_cristal WHERE id_venta = $id_venta");


//echo " id $id<br>";
// echo "id cliente $id_cliente <br>";
// echo "id cristal $id_cristal<br>";
// echo "id venta $id_venta<br>";



 if($update){
     $result = mysqli_affected_rows($conexion);
     if($result > 0){
        echo "<br><br><br><div class='row justify-content-center'><div class='alert alert-success w-20'><div class='col-md-12 text-center'>ID AGREGADO, VER PDF</div></div></div>";
        echo "  <div class='row justify-content-center'>
                    <a href='pdf/generar.php?cl=$id_cliente&v=$id' target='_blank' class='btn btn-danger'><i class='fas fa-file-pdf'></i></a>
                <div>";    
    }
        else{
            echo "<br><br><br><div class='row justify-content-center'><div class='alert alert-danger w-30'><div class='col-md-12 text-center'>VENTA INEXISTENTE</div></div></div>";
        }
    }


?>
