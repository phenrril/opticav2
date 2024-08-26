<?php
global $conexion;
include "../conexion.php";
$dniCliente = $_GET['dni'];
$query = mysqli_query($conexion, "SELECT * FROM historia_clinica WHERE dni = '$dniCliente'");

if (mysqli_num_rows($query) > 0) {
    while ($data = mysqli_fetch_assoc($query)) {
        echo '<tr>';
        echo '<td>' . $data['fecha'] . '</td>';
        echo '<td>' . $data['od_l_1'] . ' | ' . $data['od_l_2'] . ' | ' . $data['od_l_3'] . '</td>';
        echo '<td>' . $data['od_c_1'] . ' | ' . $data['od_c_2'] . ' | ' . $data['od_c_3'] . '</td>';
        echo '<td>' . $data['oi_l_1'] . ' | ' . $data['oi_l_2'] . ' | ' . $data['oi_l_3'] . '</td>';
        echo '<td>' . $data['oi_c_1'] . ' | ' . $data['oi_c_2'] . ' | ' . $data['oi_c_3'] . '</td>';
        echo '<td>' . $data['addg'] . '</td>';
        echo '<td>' . $data['armazon'] . '</td>';
        echo '<td>' . $data['precio'] . '</td>';
        echo '<td>' . $data['observaciones'] . '</td>';
        echo '</tr>';
    }
} else {
    echo '<tr><td colspan="10">No se encontró la historia clínica para el cliente con DNI ' . $dniCliente . '</td></tr>';
}

mysqli_close($conexion);
?>